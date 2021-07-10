# -*- mode: Python -*-

# Use Bazel to generate the Kubernetes YAML
watch_file('./deployments/deployment.yaml')
k8s_yaml(local('bazel run //deployments:k8s-dev'))

# Use Bazel to build the image

# The go_image BUILD rule
image_target='//cmd:image'

# Where go_image puts the image in Docker (bazel/path/to/target:name)
bazel_image='bazel/cmd:image'

custom_build(
  ref='leoh0/petstore-image',
  command=(
    'bazel run {image_target} -- --norun && ' +
    'docker tag {bazel_image} $EXPECTED_REF').format(image_target=image_target, bazel_image=bazel_image),
  deps=['./api', './cmd', './deployments'],
)

k8s_resource('petstore', port_forwards=8080)
