Feature: AzureRedisInstance feature

  @azure @allShoots @allEnvs
  Scenario: AzureRedisInstance scenario
    Given resource declaration:
      | redis      | AzureRedisInstance | "redis-"+rndStr(8)       | namespace |
      | authSecret | Secret             | redis.metadata.name      | namespace |
      | pod        | Pod                | "redis-probe-"+rndStr(8) | namespace |
    When resource redis is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AzureRedisInstance
      spec:
        redisVersion: P1
        redisConfiguration:
          maxclients: "8"
        redisVersion: "6.0"
      """

    Then eventually value load("redis").status.state equals "Ready" with timeout3X

    When resource pod is applied:
      """
      apiVersion: v1
      kind: Pod
      spec:
        containers:
        - name: redis-cli
          image: redis:latest
          command: ["/bin/bash", "-c", "--"]
          args:
          - |
            apt-get update && \
            apt-get install -y ca-certificates && \
            update-ca-certificates && \
            redis-cli -h $HOST -p $PORT -a $AUTH_STRING --tls PING
          env:
          - name: HOST
            valueFrom:
              secretKeyRef:
                key: host
                name: <(redis.metadata.name)>
          - name: PORT
            valueFrom:
              secretKeyRef:
                key: port
                name: <(redis.metadata.name)>
          - name: AUTH_STRING
            valueFrom:
              secretKeyRef:
                key: authString
                name: <(redis.metadata.name)>
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/PONG/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource redis is deleted
    Then eventually resource authSecret does not exist
    And eventually resource redis does not exist with timeout3X
