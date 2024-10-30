Feature: GcpNfsVolume feature

  @gcp @allShoots @allEnvs
  Scenario: GcpNfsVolume/Backup/Restore scenario
    Given resource declaration:
      | vol         | GcpNfsVolume          | "vol-200-gcpnfsvolume-"+rndStr(8)       | namespace |
      | pv          | PersistentVolume      | vol.status.id                           |           |
      | pvc         | PersistentVolumeClaim | vol.metadata.name                       | namespace |
      | pod         | Pod                   | "test-vol-200-gcpnfsvolume-"+rndStr(8)  | namespace |
      | backup      | GcpNfsVolumeBackup    | "backup-200-gcpnfsvolume-"+rndStr(8)    | namespace |
      | restore     | GcpNfsVolumeRestore   | "restore-200-gcpnfsvolume-"+rndStr(8)   | namespace |
      | schedule    | GcpNfsBackupSchedule  | "e2e-test-schedule-gcp-nfs-"+rndStr(8)  | namespace |
      | sch-backup  | GcpNfsVolumeBackup    | schedule.status.lastCreatedBackup.name  | namespace |
    When resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolume
      spec:
        capacityGb: 1024
      """
    Then eventually value load("vol").status.state equals "Ready" with timeout2X
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
      kind: GcpNfsVolumeBackup
      spec:
        source:
          volume:
            name: <(vol.metadata.name)>
      """
    Then eventually value load("backup").status.state equals "Ready" with timeout2X

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
              - "rm /mnt/data1/test.txt & echo 'test.txt was deleted'"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/test.txt was deleted/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource restore is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolumeRestore
      spec:
        source:
          backup:
            name: <(backup.metadata.name)>
            namespace: <(backup.metadata.namespace)>
        destination:
          volume:
            name: <(vol.metadata.name)>
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
              - "cat /mnt/data1/test.txt"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/test line/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource restore is deleted
    Then eventually resource restore does not exist
    When resource vol is deleted
    Then eventually resource pvc does not exist
    And eventually resource pv does not exist
    And eventually resource vol does not exist

    When resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolume
      spec:
        capacityGb: 1024
        sourceBackup:
          name: <(backup.metadata.name)>
          namespace: <(backup.metadata.namespace)>
      """
    Then eventually value load("vol").status.state equals "Ready" with timeout2X

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
              - "cat /mnt/data1/test.txt"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/test line/) > -1 equals true

    When resource schedule is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsBackupSchedule
      spec:
         nfsVolumeRef:
           name: <(vol.metadata.name)>
         schedule: "*/5 * * * *"
         prefix: e2e-test-backup-every-5min
         deleteCascade: true
      """
    Then eventually value load("schedule").status.state equals "Active"
    And eventually value load("schedule").status.lastCreatedBackup.name is not zero
    And eventually value load("sch-backup").status.state equals "Ready"

    When resource schedule is deleted
    Then eventually resource sch-backup does not exist
    And eventually resource schedule does not exist

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource vol is deleted
    Then eventually resource vol does not exist

    When resource backup is deleted
    Then eventually resource backup does not exist


