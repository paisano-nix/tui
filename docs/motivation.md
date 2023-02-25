# Motivation

## Problem

A Paisano-based repository has a well-defined folder structure.

That folder structure is parsed and tranfromed by Paisano into flake outputs.

However, as a project grows, and so does the number of outputs, it becomes
increasingly hard to discover all of them easily.

When it can become complicated to index for Nix power-users, it is even more
difficult for occasional repository users that are not deeply familiar with Nix.

Hence, as a project grows, we increadingly face problems of discoverability that
are not adequatly solved by the Nix CLI.

## Solution

We could, of course, write a Readme and detail all the different outputs and
bespoke ways how to interact with them.

But we all know a Readme's flaws:

- people tend to not read it, especially as it grows larger
- it tends to become outdated, coincidentially also as it grows larger

## Better Solution

To solve this discoverability problem while never becoming outdated, the
Paisano TUI / CLI parses Paisano-based flake outputs and renders them in
an easily browsable and searchable terminal user interface.

Once the user knows what she's looking for, she can switch to directly
access her intended repository action with the built-in useage as a CLI.
