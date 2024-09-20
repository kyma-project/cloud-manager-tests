Feature: GcpNfsVolumeBackup feature

  @gcp @allShoots @allEnvs
  Scenario: GcpNfsVolumeBackup scenario
    Given resource declaration:
      | vol     | GcpNfsVolume            | "vol-"+rndStr(8)        | namespace |
      | backup  | GcpNfsVolumeBackup      | "backup-"+rndStr(8)     | namespace |
    When resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolume
      spec:
        capacityGb: 1024
      """
    Then eventually value load("vol").status.state equals "Ready" with timeout2X

    When resource backup is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolumeBackup
      spec:
        source:
          volume:
            name: <(vol.metadata.name)>
      """
    Then eventually value load("backup").status.state equals "Ready" with timeout2X

    When resource backup is deleted
    Then eventually resource backup does not exist

    When resource vol is deleted
    Then eventually resource vol does not exist

