# Architecture

## Introduction
Often frameworks lock us in instead of providing us with a bootstrap to get started quickly and develop rapidly, we end
up needing to learn the framework, how it works, follow its structure guideline and if we hit a wall then we need to
dig deep into the framework, the framework needs to be updated and maintained and becomes a huge bag of things that
we constantly have to carry with us and tend to, this becomes burden over time. Golang already provides us with a 
powerful standard library that is constantly updated and which we can utilize to pretty much any extent our needs
might require.

## Hexagonal architecture

Hexagonal architecture is an approach in which we extract code/logic which communicates with "outside world" of our 
application, examples of this are database calls, network calls like http rpc sftp, UI, etc, everything that is not our
application itself. Main reasons for this is decoupling that eases maintenance and testability through mocks.
