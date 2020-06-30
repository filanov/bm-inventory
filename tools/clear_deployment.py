import os
import utils
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("-n", "--namespace", type=str, default='assisted-installer')
parser.add_argument("--delete-namespace", type=lambda x: (str(x).lower() == 'true'), default=True)
args = parser.parse_args()


def main():
    print(utils.check_output(f"kubectl delete all --all -n {args.namespace} 1> /dev/null ; true"))
    if args.delete_namespace is True:
        print(utils.check_output(f"kubectl delete namespace {args.namespace} 1> /dev/null ; true"))

if __name__ == "__main__":
    main()
