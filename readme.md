# tfc-badge
CI status badges for [Terraform Cloud](https://www.terraform.io/cloud).

## Overview

The system is quite simple:
- On each run, Terraform Cloud informs this server ("_tfcbd_") using a [webhook notification](https://www.terraform.io/docs/cloud/workspaces/notifications.html)
- _tfcbd_ stores the result in memory
- Upon request, _tfcbd_ renders a status badge that can be embedded into READMEs

## Design decisions

This project is intentionally small and simple and comes with a couple of tradeoffs:
- **No persistence.** If the server is restarted, records of previous runs are lost. 
  _tfcbd_ renders a special "unknown" badge in that case.
    
- **No scalability.** Running multiple instances of _tfcbd_ may yield wrong results.
  Honestly though, it should handle even big loads just fine and probably does not need
  to be highly available either.