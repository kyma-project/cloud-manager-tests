Feature: Main feature

  Scenario: First scenario
    Given resource declaration:
      | vol | AwsNfsVolume          |
      | pv  | PersistentVolume      |
      | pvc | PersistentVolumeClaim |
      | pod | Pod                   |
    Given resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AwsNfsVolume
      metadata:
        name: first
        namespace: default
      spec:
        capacity: 10G
      """
    Then value vol.status.state equals "True"
    Then eventually resource vol has condition Ready true
    And resource pv exists
    And resource pvc exists