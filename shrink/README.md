

# Shrink

Manages:

- Reserves new URL's
  - checks if requested URL's are already in use
  - generates randomized short URLs
    - alphanumeric characters
    - of length 5

- Retrives URLs' based off of their short code


## URL's

Maintained a hash with the following values:

short\_url ->
    LongUrl     // full length URL
    HitCount    // number of requests fulfilled


