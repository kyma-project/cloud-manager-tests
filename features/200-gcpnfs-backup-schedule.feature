Feature: GcpNfsBackupSchedule feature

  @gcp @allShoots @allEnvs
  Scenario: GcpNfsVolume scenario
    Given resource declaration:
      | vol        | GcpNfsVolume          | "vol-"+rndStr(8)                         | namespace |
      | schedule   | GcpNfsBackupSchedule  | "test-schedule-"+rndStr(3)               | namespace |
      | backup     | GcpNfsVolumeBackup    | schedule.status.lastCreatedBackup.name   | namespace |

    When resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolume
      spec:
        capacityGb: 1024
      """
    Then eventually value load("vol").status.state equals "Ready" with timeout2X

    When resource schedule is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsBackupSchedule
      spec:
         nfsVolumeRef:
           name: <(vol.metadata.name)>
         schedule: "* * * * *"
         prefix: test-minutely-backup
         deleteCascade: true
      """
    Then eventually value load("schedule").status.state equals "Active"
    And eventually value load("backup").status.state equals "Ready" with timeout2X

    When resource schedule is deleted
    Then eventually resource backup does not exist
    And eventually resource schedule does not exist

    When resource vol is deleted
    And eventually resource vol does not exist
