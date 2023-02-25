# Motivation

# Problem

A Paisano-based repository has a well-defined folder structure.

That folder structure is parsed and transformed by Paisano into flake outputs.

However, as a project grows, and so does the number of outputs, it becomes increasingly hard to discover all of them easily.

When it already becomes complicated to find stuff for Nix experts, it is even more difficult for a user who is not deeply familiar with Nix.

Hence, as a project grows, we face unsolved problems of discoverability.

# Solution

We could, of course, write a readme and detail all outputs and bespoke ways of how to interact with them.

But we all know a readme's flaws:

- people tend to not read it, especially as it grows larger
- it tends to become outdated, coincidentally also as it grows larger

# Better Solution

To solve this discoverability problem while never becoming outdated, this tool renders all Paisano-based flake outputs into an easily browsable and searchable terminal user interface. A variety of actions offer the user pre-defined and discoverable ways to interact with these outputs based on their types.

Once the user knows what she's looking for, she can choose to access her intended target action directly by using this tool as a CLI.
