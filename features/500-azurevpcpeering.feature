Feature: AzureVpcPeering feature

  @azure @allShoots @allEnvs
  Scenario: AzureVpcPeering scenario
    Given resource declaration:
      | peering | AzureVpcPeering | "peering-"+rndStr(8)       | namespace |
      | pod     | Pod             | "peering-probe-"+rndStr(8) | namespace |
    When resource peering is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AzureVpcPeering
      spec:
        remotePeeringName: e2e
        remoteVnet: >-
          /subscriptions/3f1d2fbd-117a-4742-8bde-6edbcdee6a04/resourceGroups/e2e/providers/Microsoft.Network/virtualNetworks/e2e
        deleteRemotePeering: true
      """

    Then eventually value load("peering").status.state equals "Connected"

    When resource pod is applied:
      """
      apiVersion: v1
      kind: Pod
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
            image: alpine
            command:
              - "nc"
            args:
              - "-zv"
              - "172.23.0.4"
              - "22"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/172.23.0.4 \(172.23.0.4:22\) open/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource peering is deleted
    Then eventually resource peering does not exist with timeout3X
