Feature: GcpNfsVolumeRestore-NewVolume feature

  @gcp @allShoots @allEnvs
  Scenario: GcpNfsVolumeRestore-NewVolume scenario
    Given resource declaration:
      | vol     | GcpNfsVolume            | "vol-"+rndStr(8)        | namespace |
      | backup  | GcpNfsVolumeBackup      | "backup-"+rndStr(8)     | namespace |
      | vol2    | GcpNfsVolume            | "vol-"+rndStr(8)        | namespace |
      | pod     | Pod                     | "test-vol"              | namespace |
      | pod2    | Pod                     | "test-vol2"             | namespace |

    When resource vol is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolume
      spec:
        capacityGb: 1024
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
              - "echo 'test line' > /mnt/data1/test.txt & cat /mnt/data1/test.txt"
        restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/test line/) > -1 equals true


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

    When resource vol2 is applied:
      """
      apiVersion: cloud-resources.kyma-project.io/v1beta1
      kind: GcpNfsVolume
      spec:
        capacityGb: 1024
        sourceBackup:
          name: <(backup.metadata.name)>
          namespace: <(backup.metadata.namespace)>
      """
    Then eventually value load("vol2").status.state equals "Ready" with timeout2X

    When resource pod2 is applied:
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
    Then eventually value load("pod2").status.phase equals "Succeeded"
    And value logs("pod2").search(/test line/) > -1 equals true

    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource pod2 is deleted
    Then eventually resource pod2 does not exist

    When resource vol2 is deleted
    Then eventually resource vol2 does not exist

    When resource backup is deleted
    Then eventually resource backup does not exist

    When resource vol is deleted
    Then eventually resource vol does not exist



