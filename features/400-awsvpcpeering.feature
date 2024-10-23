Feature: AwsVpcPeering feature

  @aws @allShoots @allEnvs
  Scenario: AwsVpcPeering scenario
    Given resource declaration:
      | peering | AwsVpcPeering | "peering-"+rndStr(8)       | namespace |
      | pod     | Pod           | "peering-probe-"+rndStr(8) | namespace |
    When resource peering is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AwsVpcPeering
      metadata:
        name: e2e
      spec:
        remoteAccountId: "642531956841"
        remoteRegion: "us-east-1"
        remoteVpcId: "vpc-0709fb45c2be50920"
        deleteRemotePeering: true
      """

    Then eventually value load("peering").status.state equals "Ready" with

    When resource pod is applied:
      """
      apiVersion: v1
      kind: Pod
      metadata:
        name: awsvpcpeering-demo
      spec:
        containers:
          - name: my-container
            resources:
              limits:
                memory: 512Mi
                cpu: "1"
              requests:
                memory: 256Mi
                cpu: "0.2"
            image: ubuntu
            command:
              - "/bin/bash"
              - "-c"
              - "--"
            args:
              - "apt update; apt install netcat-openbsd -y; nc -zv 10.3.136.98"
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/Connection to 10.3.136.98 22 port \[tcp\/\*\] succeeded!/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource peering is deleted
    Then eventually resource peering does not exist with timeout3X
