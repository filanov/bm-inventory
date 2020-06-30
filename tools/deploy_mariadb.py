import os
import utils
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("-n", "--namespace", help='namespace to use', type=str, default='assisted-installer')
args = parser.parse_args()


def main():
    src_file = os.path.join(os.getcwd(), "deploy/mariadb/mariadb-configmap.yaml")
    dst_file = os.path.join(os.getcwd(), "build/mariadb-configmap.yaml")
    with open(src_file, "r") as src:
        with open(dst_file, "w+") as dst:
            data = src.read()
            print("Deploying {}".format(dst_file))
            dst.write(data)

    utils.apply(dst_file)

    src_file = os.path.join(os.getcwd(), "deploy/mariadb/mariadb-deployment.yaml")
    dst_file = os.path.join(os.getcwd(), "build/mariadb-deployment.yaml")
    with open(src_file, "r") as src:
        with open(dst_file, "w+") as dst:
            data = src.read()
            print("Deploying {}".format(dst_file))
            dst.write(data)
    utils.apply(dst_file)

    src_file = os.path.join(os.getcwd(), "deploy/mariadb/mariadb-storage.yaml")
    dst_file = os.path.join(os.getcwd(), "build/mariadb-storage.yaml")
    with open(src_file, "r") as src:
        with open(dst_file, "w+") as dst:
            data = src.read()
            try:
                size = utils.check_output(
                    f"kubectl -n {args.namespace} get persistentvolumeclaims mariadb-pv-claim " +
                    "-o=jsonpath='{.status.capacity.storage}'")
                print("Using existing disk size", size)
            except:
                size = "10Gi"
                print("Using default size", size)
            data = data.replace("REPLACE_STORAGE", size)
            print("Deploying {}".format(dst_file))
            dst.write(data)

    utils.apply(dst_file)


if __name__ == "__main__":
    main()
