#!/bin/bash

#
# SPDX-License-Identifier: BSD-3-Clause
# Copyright(©) 2024 Intel Corporation
# Intel® Tiber™ Broadcast Suite
#

. VARIABLES.rc

# Check if VFIO_PORT_T is set
if [ -z "$VFIO_PORT_T" ]; then
    echo -e "\e[31mError: VFIO_PORT_T is not set.\e[0m"
    echo "Use dpdk-devbind.py -s to check pci address of vfio device"
    exit 1
fi

docker run -it \
   --user root\
   --privileged \
   --device=/dev/vfio:/dev/vfio \
   --device=/dev/dri:/dev/dri \
   --cap-add ALL \
   -v "$(pwd)":/videos \
   -v /usr/lib/x86_64-linux-gnu/dri:/usr/local/lib/x86_64-linux-gnu/dri/ \
   -v /tmp/kahawai_lcore.lock:/tmp/kahawai_lcore.lock \
   -v /dev/null:/dev/null \
   -v /tmp/hugepages:/tmp/hugepages \
   -v /hugepages:/hugepages \
   -v /var/run/imtl:/var/run/imtl \
   --network=my_net_801f0 \
   --ip=192.168.2.1 \
   --expose=20000-20170 \
   --ipc=host -v /dev/shm:/dev/shm \
      video_production_image \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_1.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_2.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_1.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_2.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_1.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_2.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_1.yuv \
        -video_size 1920x1080 -f rawvideo -pix_fmt yuv422p10le -i /videos/src/1080p_yuv422_10b_2.yuv \
        -map 0:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20000 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 1:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20001 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 2:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20002 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 3:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20003 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 4:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20004 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 5:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20005 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 6:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20006 -payload_type 112 -fps 25 -f mtl_st20p - \
        -map 7:v -p_port "${VFIO_PORT_T}" -p_sip 192.168.2.1 -p_tx_ip 192.168.2.2 -udp_port 20007 -payload_type 112 -fps 25 -f mtl_st20p -
