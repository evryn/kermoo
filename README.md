[![integration tests](https://github.com/evryn/kermoo/actions/workflows/ci-build.yaml/badge.svg?branch=main)](https://github.com/evryn/kermoo/actions/workflows/ci-build.yaml) [![codecov](https://codecov.io/gh/evryn/kermoo/branch/main/graph/badge.svg)](https://codecov.io/gh/evryn/kermoo) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=evryn_kermoo&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=evryn_kermoo)  

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=evryn_kermoo&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=evryn_kermoo) 
[![CodeQL](https://github.com/evryn/kermoo/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/evryn/kermoo/actions/workflows/github-code-scanning/codeql)

<p align="center">
<img src="./docs/kermoo.png" width="200" height="200">
</p>
<h1 align="center">Kermoo</h1>
<p align="center">
🪱 A Delightfully Buggy Application 🪱
</p>


---

<center>
<p><strong>⚠️ Kermoo is heavily under development and doesn't have a stable release yet. Keep checking this page again in a day or two. ⚠️</strong></p>
</center>

---

## Introduction

Ahoy there, brave DevOps and SRE warriors! 🏴‍☠️ Introducing *Kermoo*, derived from the Persian mischief-maker "کرمو". We're here to toss a wrench (or ten) into your perfectly running apps, all in the name of science... and a bit of fun. 😜

Want to see how your system fares against the worst? Or perhaps you've got a tad too much trust in your resilience and fault tolerance? Time to put that to the test!

## 🚀 Features

1. **🕐 Simulate Process Startup Delays & Failures**:
    - Make your process sleep in a bit. They deserve it! Define a delay.
    - Surprise! Unexpected exits. 🎉 Set an exit time and code.

2. **💥 Simulate Webserver & Backend Mayhem**:
    - Chaos in the form of HTTP requests: 
        * Oops! Request failed. 🙈
        * Snail-paced 🐌 responses.
        * Lost at connection level! 📞❌
    - Plan your mischief: percentage affected, duration...
    - Or sometimes...just sometimes, send some good [static or whoami-like] response. 🌈

3. **🔥 Simulate Heavy CPU Sunbathing**:
    - Turn up the heat and get that CPU sweating! 💦 Set a load percentage and duration.

4. **🧠 Simulate Forgetful Memory Leaks**:
    - Because who doesn't want to spring a leak now and then? Choose your memory size and duration.

## 🔆 Installation
Kermoo is available as:
- Docker Image
- Helm Chart
- Downloadable Binaries

<center>Read More: <strong><a href="https://github.com/evryn/kermoo/wiki">Installation and Getting Started</a></strong></center>

## 🛠 Setup & Configuration

Simply setup the all-optional configurations using YAML and pass it to Kermoo:

```yaml
kermoo start -f - <<EOL
  # Simulate the main process to delay initially for 5 seconds
  # Then simulate disaster by exiting the process after somewhere
  # between 4 to 10 seconds with the exit code of 20.
  process:
    delay: 5s
    exit:
      after: 4s to 10s
      code: 20

  webServers:
    # Setup a webserver that listens on 0.0.0.0:80
    # that will be ready for connect 50% of time
    - fault:
        percentage: 50
      routes:
        # Setup a /livez route that will serve a whoami content
        # with the chance of 60% or fail with an error (4xx, 5xx)
        - path: /livez
          content:
            whoami: true
          fault:
            percentage: 60

  # Simulate CPU load which will repeatitivly utilize
  # 20%, 50% and 80% of all cores every second.
  cpuLoad:
    percentage: 20, 50, 80
    interval: 1s

  # Simulate memory leak which will repeatitivly consume
  # somewhere between 100Mi and 1Gi of memory - recomputed
  # every 5 seconds.
  memoryLeak:
    size: 100Mi to 1Gi
    interval: 5s
EOL
```

<center>Read More: <a href="https://github.com/evryn/kermoo/wiki">Advanced Configuration</a></center>

## 🤝 Join the Chaos Club
Love the pandemonium? Consider [contributing](CONTRIBUTING.md) and let's brew more bedlam together! ☕🔥
