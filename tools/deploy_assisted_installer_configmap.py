# -*- coding: utf-8 -*-
import os
import utils
import argparse
import yaml
import deployment_options


SRC_FILE = os.path.join(os.getcwd(), "deploy/bm-inventory-configmap.yaml")
DST_FILE = os.path.join(os.getcwd(), "build/bm-inventory-configmap.yaml")
SERVICE = "bm-inventory"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--target")
    parser.add_argument("--domain")
    deploy_options = deployment_options.load_deployment_options(parser)
    # TODO: delete once rename everything to assisted-installer
    if deploy_options.target == "oc-ingress":
        service_host = "assisted-installer.{}".format(
            utils.get_domain(deploy_options.domain)
        )
        service_port = "80"
    else:
        service_host = utils.get_service_host(SERVICE, deploy_options.target)
        service_port = utils.get_service_port(SERVICE, deploy_options.target)
    with open(SRC_FILE, "r") as src:
        with open(DST_FILE, "w+") as dst:
            data = src.read()
            data = data.replace("REPLACE_URL", '"{}"'.format(service_host))
            data = data.replace("REPLACE_PORT", '"{}"'.format(service_port))
            print("Deploying {}".format(DST_FILE))

            versions = {
                "IMAGE_BUILDER": "installer-image-build",
                "AGENT_DOCKER_IMAGE": "agent",
                "KUBECONFIG_GENERATE_IMAGE": "ignition-manifests-and-kubeconfig-generate",
                "INSTALLER_IMAGE": "assisted-installer",
                "CONNECTIVITY_CHECK_IMAGE": "connectivity_check",
                "INVENTORY_IMAGE": "inventory",
                "HARDWARE_INFO_IMAGE": "hardware_info",
            }
            for env_var_name, image_short_name in versions.items():
                image_fqdn = deployment_options.get_image_override(
                    deploy_options, image_short_name, env_var_name
                )
                versions[env_var_name] = image_fqdn
            versions["SELF_VERSION"] = deployment_options.get_image_override(
                deploy_options, "bm-inventory", "SERVICE"
            )
            y = yaml.load(data)
            y["data"].update(versions)
            data = yaml.dump(y)
            dst.write(data)

    utils.apply(DST_FILE)


if __name__ == "__main__":
    main()
