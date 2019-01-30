#!/usr/bin/env bash

set -xeu

os=$1
arch=$2

seccomp_opts() {
  local docker_ver minor_ver
  docker_ver=$(docker version -f '{{.Client.Version}}')
  minor_ver=$(cut -d'.' -f2 <<<$docker_ver)
  if ((minor_ver>=10)); then
    echo '--security-opt seccomp:unconfined'
  fi
}

cp "../release/goss-linux-$arch" "goss/$os/"
# Run build if Dockerfile has changed but hasn't been pushed to dockerhub
if ! md5sum -c "Dockerfile_${os}.md5"; then
  docker build -q -t "simonbaeumer/goss_${os}:latest" - < "Dockerfile_$os"
# Pull if image doesn't exist locally
elif ! docker images | grep "SimonBaeumer/goss_$os";then
  docker pull "simonbaeumer/goss_$os"
fi

container_name="goss_int_test_${os}_${arch}"
docker_exec() {
  docker exec "$container_name" "$@"
}

# Cleanup any old containers
if docker ps -a | grep "$container_name";then
  docker rm -vf "$container_name"
fi
opts=(--env OS=$os --cap-add SYS_ADMIN -v "$PWD/goss:/goss"  -d --name "$container_name" $(seccomp_opts))
id=$(docker run "${opts[@]}" "simonbaeumer/goss_$os" /sbin/init)
ip=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' "$id")
trap "rv=\$?; docker rm -vf $id; exit \$rv" INT TERM EXIT
# Give httpd time to start up, adding 1 second to see if it helps with intermittent CI failures
[[ $os != "arch" ]] && docker_exec "/goss/$os/goss-linux-$arch" -g "/goss/goss-wait.yaml" validate -r 10s -s 100ms && sleep 1

#out=$(docker exec "$container_name" bash -c "time /goss/$os/goss-linux-$arch -g /goss/$os/goss.yaml validate")
out=$(docker_exec "/goss/$os/goss-linux-$arch" --vars "/goss/vars.yaml" -g "/goss/$os/goss.yaml" validate)
echo "$out"

if [[ $os == "arch" ]]; then
  egrep -q 'Count: 74, Failed: 0' <<<"$out"
else
  egrep -q 'Count: 88, Failed: 0' <<<"$out"
fi

if [[ ! $os == "arch" ]]; then
  docker_exec /goss/generate_goss.sh "$os" "$arch"

  #docker exec $container_name bash -c "cp /goss/${os}/goss-generated-$arch.yaml /goss/${os}/goss-expected.yaml"
  docker_exec diff -wu "/goss/${os}/goss-expected.yaml" "/goss/${os}/goss-generated-$arch.yaml"

  #docker exec $container_name bash -c "cp /goss/${os}/goss-aa-generated-$arch.yaml /goss/${os}/goss-aa-expected.yaml"
  docker_exec diff -wu "/goss/${os}/goss-aa-expected.yaml" "/goss/${os}/goss-aa-generated-$arch.yaml"

  docker_exec /goss/generate_goss.sh "$os" "$arch" -q

  #docker exec $container_name bash -c "cp /goss/${os}/goss-generated-$arch.yaml /goss/${os}/goss-expected-q.yaml"
  docker_exec diff -wu "/goss/${os}/goss-expected-q.yaml" "/goss/${os}/goss-generated-$arch.yaml"
fi

#docker rm -vf goss_int_test_$os
