Feature: GcpVpcPeering feature
  @gcp @allShoots @allEnvs
  Scenario: GcpVpcPeering scenario
    Given resource declaration:
      | vpcPeering | GcpVpcPeering | "vpcpeering-"+rndStr(8)          | namespace |
      | pod        | Pod           | "vpcpeering-test-pod-"+rndStr(8) | namespace |
    When resource vpcPeering is applied:
        """
        apiVersion: cloud-resources.kyma-project.io/v1beta1
        kind: GcpVpcPeering
        metadata:
          name: "gcp-vpc-peering-e2e-test"
        spec:
          remotePeeringName: "vpc-peering-e2e-tests-to-sap-gcp-skr-dev-cust-00002"
          remoteProject: "sap-sc-learn"
          remoteVpc: "vpc-peering-e2e-tests"
          importCustomRoutes: false
        """
    Then eventually value load("vpcPeering").status.state equals "Connected"

    When resource pod is applied:
        """
        apiVersion: v1
        kind: Pod
        spec:
          containers:
          - name: netcat
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
              - "10.240.254.2"
              - "22"
          restartPolicy: Never
      """
    Then eventually value load("pod").status.phase equals "Succeeded"
    And value logs("pod").search(/10.240.254.2 \(10.240.254.2:22\) open/) > -1 equals true


    When resource pod is deleted
    Then eventually resource pod does not exist

    When resource vpcPeering is deleted
    Then eventually resource vpcPeering does not exist