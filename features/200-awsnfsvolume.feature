Feature: AwsNfsVolume feature

  @aws @allShoots @allEnvs
  Scenario: AwsNfsVolume scenario
    Given resource declaration:
      | vol         | AwsNfsVolume          | "vol-200-awsnfsvolume-"+rndStr(8)       | namespace |
      | pv          | PersistentVolume      | vol.status.id                           |           |
      | pvc         | PersistentVolumeClaim | vol.metadata.name                       | namespace |
      | pod         | Pod                   | "test-vol-200-awsnfsvolume-"+rndStr(8)  | namespace |
      | backup      | AwsNfsVolumeBackup    | "backup-200-awsnfsvolume-"+rndStr(8)    | namespace |
      | restore     | AwsNfsVolumeRestore   | "restore-200-awsnfsvolume-"+rndStr(8)   | namespace |
      | schedule    | AwsNfsBackupSchedule  | "e2e-test-schedule-aws-nfs-"+rndStr(8)  | namespace |
      | sch-backup  | AwsNfsVolumeBackup    | schedule.status.lastCreatedBackup.name  | namespace |
    When resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AwsNfsVolume
      spec:
        capacity: 10G
      """
    Then eventually value load("vol").status.state equals "Ready"
    And eventually value load("pv").status.phase equals "Bound"
    And eventually value load("pvc").status.phase equals "Bound"

    When resource pod is applied:
      """
      apiVersion: v1
      kind: Pod
      spec:
        volumes:
          - name: data
            persistentVolumeClaim:
              claimName: <(vol.metadata.name)>
        containers:
          - name: cloud1
            image: ubuntu
            imagePullPolicy: IfNotPresent
            volumeMounts:
              - mountPath: "/mnt/data1"
                name: data
            command:
              - "/bin/bash"
              - "-c"
              - "--"
            args:
              - "echo 'test line' > /mnt/data1/test.txt & cat /mnt/data1/test.txt"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/test line/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource backup is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AwsNfsVolumeBackup
      spec:
        source:
          volume:
            name: <(vol.metadata.name)>
      """
    Then eventually value load("backup").status.state equals "Ready" with timeout2X

    When resource restore is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: AwsNfsVolumeRestore
      spec:
        source:
          backup:
            name: <(backup.metadata.name)>
            namespace: <(backup.metadata.namespace)>
      """
    Then eventually value load("restore").status.state equals "Done" with timeout2X

    When resource pod is applied:
      """
      apiVersion: v1
      kind: Pod
      spec:
        volumes:
          - name: data
            persistentVolumeClaim:
              claimName: <(vol.metadata.name)>
        containers:
          - name: cloud1
            image: ubuntu
            imagePullPolicy: IfNotPresent
            volumeMounts:
              - mountPath: "/mnt/data1"
                name: data
            command:
              - "/bin/bash"
              - "-c"
              - "--"
            args:
              - "cat /mnt/data1/aws-backup-restore*/test.txt"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/test line/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource vol is deleted
    Then eventually resource pvc does not exist
    And eventually resource pv does not exist
    And eventually resource vol does not exist

    When resource backup is deleted
    Then eventually resource backup does not exist

    When resource restore is deleted
    Then eventually resource restore does not exist
