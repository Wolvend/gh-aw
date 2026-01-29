---
name: Custom Binary Mounts Demo
description: Demonstrates using custom binary and library mounts with AWF sandbox
engine: copilot
on:
  workflow_dispatch:

# Configure sandbox with custom tool mounts
sandbox:
  agent:
    id: awf
    mounts:
      # Database clients
      - "/usr/bin/psql:/usr/bin/psql:ro"
      - "/usr/bin/mysql:/usr/bin/mysql:ro"
      
      # Cloud CLIs
      - "/usr/local/bin/aws:/usr/local/bin/aws:ro"
      
      # Build tools
      - "/usr/bin/make:/usr/bin/make:ro"
      - "/usr/bin/gcc:/usr/bin/gcc:ro"
      
      # Shared libraries
      - "/usr/lib/x86_64-linux-gnu/libssl.so.3:/usr/lib/x86_64-linux-gnu/libssl.so.3:ro"

network:
  firewall: true
  allowed:
    - defaults

tools:
  github:
    mode: remote
---

# Custom Binary Mounts Demonstration

This workflow demonstrates how to use custom binary and library mounts to make
specialized tools available inside the AWF sandbox container.

## Task

Please verify that the following custom-mounted binaries are available:

1. Check for database clients:
   - Run `which psql` to verify PostgreSQL client is available
   - Run `which mysql` to verify MySQL client is available

2. Check for cloud CLIs:
   - Run `which aws` to verify AWS CLI is available

3. Check for build tools:
   - Run `which make` to verify Make is available
   - Run `which gcc` to verify GCC compiler is available

4. Check for shared libraries:
   - Run `ldconfig -p | grep libssl` to verify SSL library is available

Display the results in a summary showing which tools are available and which
are missing. If a tool is available, also show its version.
