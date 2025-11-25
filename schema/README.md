# Schema Definition
Definition and logic behind the schema for `casting.yml` - [reference](https://www.notion.so/signoz/Foundry-263fcc6bcd19804c89efcbb7bae9effb?source=copy_link#26bfcc6bcd198014b184de990757ff04)
- The schema is subject to change as it is in early development
- Schema is used by `foundryctl` to generate the moldings for each platform

# Requirements
- Cue installed `brew install cue-lang/tap/cue`
- Text editor with cue lang highlights (VScode Extension available)


# Usage
- Navigate to the each platform directory
- Use cue to validate the schema is correct `cue vet -c schema/casting-schema.cue casting.yml`
- Cue will outline if your `casting.yml` matches the schema
