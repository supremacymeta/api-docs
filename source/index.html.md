---
title: Supremacy API Reference

language_tabs: # must be one of https://git.io/vQNgJ
  - javascript

search: true

toc_footers:
  - <a href="https://supremacy.game">Supremacy</a>
  - <a href="https://play.supremacy.game">Battle Arena</a>
  - <a href="https://supremacygame.dev">Proving Grounds</a>
code_clipboard: true

meta:
  - name: description
    content: Documentation for the Supremacy API
---

# Introduction

Welcome to the Supremacy API.

This document covers API documentation and examples for you to develop applications and analytics for the Battle Arena, and eventually the rest of the system.

The response request and response types can change at any moment.

# Environments

- Production: `https://api.supremacy.game/api`
- Proving Grounds: `https://api.supremacygame.dev/api`

# Authentication

No authentication is required for the public API endpoints.

# Signature Verification

The verified API endpoints provided with Supremacy are signed by the contract operator:

- Staging: `0xc01c2f6DD7cCd2B9F8DB9aa1Da9933edaBc5079E`
- Production: `0xeCfB1f31F012Db0bf6720610301F23F064c567f9`

Developers can use this to commit verified data into the blockchain in a trustless manner.

## Signature Verification

```solidity
// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

/// @custom:security-contact privacy-admin@supremacy.game
contract SignatureVerifier {
    address public signer;

    constructor(address _signer) {
        require(_signer != address(0), "zero address can not be signer");
        signer = _signer;
    }

    function setSigner(address _signer) internal {
        require(_signer != address(0), "zero address can not be signer");
        signer = _signer;
        emit SetSigner(_signer);
    }

    // verify returns true if signature by signer matches the hash
    function verify(bytes32 messageHash, bytes memory signature)
        internal
        view
        returns (bool)
    {
        require(signer != address(0), "zero address can not be signer");
        bytes32 ethSignedMessageHash = getEthSignedMessageHash(messageHash);
        (bytes32 r, bytes32 s, uint8 v) = splitSignature(signature);
        return recoverSigner(ethSignedMessageHash, r, s, v) == signer;
    }

    // verifyRSV returns true if signature by signer matches the hash
    function verifyRSV(
        bytes32 messageHash,
        bytes32 r,
        bytes32 s,
        uint8 v
    ) internal view returns (bool) {
        require(signer != address(0), "zero address can not be signer");
        bytes32 ethSignedMessageHash = getEthSignedMessageHash(messageHash);
        return recoverSigner(ethSignedMessageHash, r, s, v) == signer;
    }

    function getEthSignedMessageHash(bytes32 messageHash)
        internal
        pure
        returns (bytes32)
    {
        return
            keccak256(
                abi.encodePacked(
                    "\x19Ethereum Signed Message:\n32",
                    messageHash
                )
            );
    }

    function recoverSigner(
        bytes32 _ethSignedMessageHash,
        bytes32 r,
        bytes32 s,
        uint8 v
    ) internal pure returns (address) {
        require(v == 27 || v == 28, "invalid v value");
        require(
            uint256(s) <
                0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A1,
            "invalid s value"
        );
        return ecrecover(_ethSignedMessageHash, v, r, s);
    }

    function splitSignature(bytes memory signature)
        internal
        pure
        returns (
            bytes32 r,
            bytes32 s,
            uint8 v
        )
    {
        require(signature.length == 65, "invalid signature length");
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
        }
    }

    event SetSigner(address _signer);
}


```

Use [EIP-1271](https://eips.ethereum.org/EIPS/eip-1271) to verify signatures:

- Split signature from API into `r`, `s`, `v` components
- Generate message hash using the arguments used to generate the hash (see below)
- Use `ecrecover(_ethSignedMessageHash, v, r, s);` to retrieve the signer
- Confirm signer matches the correct public address
  - Staging: `0xc01c2f6DD7cCd2B9F8DB9aa1Da9933edaBc5079E`
  - Production: `0xeCfB1f31F012Db0bf6720610301F23F064c567f9`

## Previous Battle Hash Generation

```solidity
// BattleCommitHash returns the hash of a historical battle record for signature verification
function BattleCommitHash(BattleCommit memory _battleCommit)
    internal
    pure
    returns (bytes32)
{
    return
        keccak256(
            abi.encode(
                battle_number,
                battle_started_at,
                battle_ended_at,
                battle_winner,
                battle_runner_up,
                battle_loser
            )
        );
}

```

Battle commits use 6 values to generate their signatures:

- `uint256 battle_number`: The battle number
- `uint256 battle_started_at`: The time in unix seconds (int64) that the battle started
- `uint256 battle_ended_at`: The time in unix seconds (int64) that the battle ended
- `uint256 battle_winner`: The battle winner
- `uint256 battle_runner_up`: The battle runner-up
- `uint256 battle_loser`: The battle loser

The factions are mapped to integers, namely:

- 1: Zaibatsu
- 2: Red Mountain
- 3: Boston Cybernetics

Retrieve these from the REST API to generate the message hash.

## Current Battle Hash Generation

```solidity
// CurrentBattleNumberHash returns the hash of the current battle for signature verification
function CurrentBattleNumberHash(CurrentBattle memory _currentBattle)
    internal
    pure
    returns (bytes32)
{
    return
        keccak256(
            abi.encode(
                battle_number,
                battle_started_at,
                battle_expires_at
            )
        );
}
```

Current battles 3 values to generate their signatures.

- `uint256 battle_number`: The battle number
- `uint256 battle_started_at`: The time in unix seconds (int64) that the battle started
- `uint256 battle_expires_at`: A short expiry (15 seconds) which invalidates the signatures to prevent repeat usage.

Retrieve these from the REST API to generate the message hash.

# Battle History

## Get Latest Battles

```javascript
fetch("http://api.supremacygame.dev/api/battle_history")
  .then((resp) => {
    return resp.json();
  })
  .then((resp) => {
    console.log(resp);
  })
  .catch((err) => {
    console.error(err);
  });
```

> The above command returns JSON structured like this:

```json
{
  "current_battle": {
    "number": 11,
    "started_at": 1662527943,
    "expires_at": 1662528117,
    "signature": "0x333c19d4ea40a8c258f7a5084d1565e9cb4145c6273c0ea7c68124f4e42d8374663a3fd873123ad2e07f726d9c6c4b93a7857c03d0959c8b0bd8f788ad1fd7a01b"
  },
  "previous_battles": [
    {
      "number": 10,
      "started_at": 1662525218,
      "ended_at": 1662525502,
      "winner": 1,
      "runner_up": 3,
      "loser": 2,
      "signature": "0xfa13ca633babd42aa0944e6213c1d657f075637de6690eb2c9dbae5c05d3e36f455bedbbc3ce4e82fe73d9756f72d56056d54092a4581712d19f16ae23f236db1c"
    },
    ...
  ]
}
```

This endpoint retrieves the current battle, and the latest ten battles.

### HTTP Request

`GET https://api.supremacygame.dev/api/battle_history`

### Query Parameters

None.

## Get a battle history record

```javascript
fetch("http://api.supremacygame.dev/api/battle_history/1")
  .then((resp) => {
    return resp.json();
  })
  .then((resp) => {
    console.log(resp);
  })
  .catch((err) => {
    console.error(err);
  });
```

> The above command returns JSON structured like this:

```json
{
  "battle": {
    "number": 1,
    "started_at": 1662518516,
    "ended_at": 1662518992,
    "winner": 2,
    "runner_up": 3,
    "loser": 1,
    "signature": "0x9b2760283db0586ffd5bca27858ec50437dc49fe73ba8b5132de71a7359cc0b6505ce167de06fa3ebd6c999725086930864c8a3599c81ea40bff3793669b708f1b"
  }
}
```

This endpoint retrieves a battle record, given its ID.

### HTTP Request

`GET https://api.supremacygame.dev/api/battle_history/{battle_number}`

### URL Parameters

| Parameter     | Description                   |
| ------------- | ----------------------------- |
| battle_number | The battle number to retrieve |
