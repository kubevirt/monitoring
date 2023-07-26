import argparse
import json
import os
import shutil


def build_ghpages():
    config_file, output_dir = parse_flags()
    parsed_content = parse_config_file(config_file)
    build_output(parsed_content, output_dir)


def parse_flags():
    parser = argparse.ArgumentParser()
    parser.add_argument("--config_file", help="Runbooks configuration file", required=True)
    parser.add_argument("--output_dir", help="Output directory", default="dist")
    args = parser.parse_args()
    return args.config_file, args.output_dir


def parse_config_file(config_file):
    with open(config_file) as f:
        content = f.read()
    return json.loads(content)


def build_output(config, output_dir):
    validate_key(config, "metrics")
    validate_key(config, "runbooks")
    validate_key(config, "deprecated_runbooks")
    validate_key(config, "redirects")

    reset_dir(output_dir)

    build_metrics(config["metrics"], output_dir)
    runbooks = build_runbooks(config["runbooks"], output_dir)
    deprecate_runbooks = build_runbooks(config["deprecated_runbooks"], output_dir)
    redirected_runbooks = link_renamed_runbooks(config["redirects"], output_dir)

    build_index(output_dir)
    build_runbooks_index(output_dir, runbooks, deprecate_runbooks, redirected_runbooks)


def reset_dir(output_dir):
    if os.path.exists(output_dir):
        shutil.rmtree(output_dir)
    os.makedirs(output_dir)


def validate_key(config, key):
    if key not in config:
        raise Exception(f"'{key}' key not found in config file")


def build_metrics(metrics, output_dir):
    if not os.path.exists(metrics):
        raise Exception(f"metrics file '{metrics}' not found")

    shutil.copy(metrics, os.path.join(output_dir, "metrics.md"))


def build_runbooks(original_dir, output_dir):
    files = []

    final_runbooks_dir = os.path.join(output_dir, "runbooks")

    if not os.path.exists(final_runbooks_dir):
        os.makedirs(final_runbooks_dir)

    for runbook in os.listdir(original_dir):
        if runbook.endswith(".md") and runbook != "README.md":
            files.append(runbook)
            shutil.copy(os.path.join(original_dir, runbook), os.path.join(final_runbooks_dir, runbook))

    return files


def link_renamed_runbooks(redirects, output_dir):
    files = []

    for redirect in redirects:
        validate_key(redirect, "source")
        validate_key(redirect, "target")

        source = os.path.join(output_dir, "runbooks", redirect["source"])
        target = os.path.join(output_dir, "runbooks", redirect["target"])
        if os.path.exists(target):
            files.append((redirect["source"], redirect["target"]))
            shutil.copy(target, source)
        else:
            print(f"WARNING: target file '{target}' not found, skipping redirect")

    return files


def build_index(output_dir):
    index = os.path.join(output_dir, "index.md")
    with open(index, "w") as f:
        f.write("## Metrics\n\n")
        f.write(f"- [Metrics](metrics.md)\n\n")
        f.write("## Runbooks\n\n")
        f.write(f"- [Runbooks](runbooks_index.md)\n\n")


def build_runbooks_index(output_dir, runbooks, deprecate_runbooks, redirected_runbooks):
    index = os.path.join(output_dir, "runbooks_index.md")
    with open(index, "w") as f:
        f.write("# KubeVirt Runbooks\n\n")

        for runbook in runbooks:
            f.write(f"- [{runbook}](runbooks/{runbook})\n")
        f.write("\n")

        f.write("## Deprecated Runbooks\n\n")
        for runbook in deprecate_runbooks:
            f.write(f"- [{runbook}](runbooks/{runbook})\n")
        f.write("\n")

        f.write("## Renamed Runbooks\n\n")
        for runbook in redirected_runbooks:
            f.write(f"- [{runbook[0]} -> {runbook[1]}](runbooks/{runbook[1]})\n")
        f.write("\n")


if __name__ == '__main__':
    build_ghpages()
