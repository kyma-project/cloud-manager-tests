Feature: GcpRedisInstance feature

  @gcp @allShoots @allEnvs
  Scenario: GcpRedisInstance scenario
    Given resource declaration:
      | redis      | GcpRedisInstance | "redis-"+rndStr(8)       | namespace |
      | authSecret | Secret           | redis.metadata.name      | namespace |
      | pod        | Pod              | "redis-probe-"+rndStr(8) | namespace |
    When resource redis is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpRedisInstance
      spec:
        memorySizeGb: 5
        tier: "STANDARD_HA"
        redisVersion: REDIS_7_1
        authEnabled: true
        transitEncryption:
          serverAuthentication: true
        redisConfigs:
          maxmemory-policy: volatile-lru
          activedefrag: "yes"
        maintenancePolicy:
          dayOfWeek:
            day: "SATURDAY"
            startTime:
                hours: 15
                minutes: 45
      """
    Then eventually value load("redis").status.state equals "Ready" with timeout5X

    When resource pod is applied:
      """
      apiVersion: v1
      kind: Pod
      spec:
        containers:
        - name: redis-cli
          image: redis:latest
          command: ["/bin/bash", "-c", "--"]
          args: ["redis-cli -h $HOST -p $PORT -a $AUTH_STRING --tls --cacert /mnt/CaCert.pem PING"]
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
          volumeMounts:
          - name: mounted
            mountPath: /mnt
        volumes:
        - name: mounted
          secret:
            secretName: <(redis.metadata.name)>
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/PONG/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource redis is deleted
    Then eventually resource authSecret does not exist
    And eventually resource redis does not exist
