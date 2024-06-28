#!/bin/bash 
cat partial_unblock_udp_fragment.json | sudo bin/fastnetmon_flow_spec_fragmentation 
