/*!
 * spvnode.js - SPV-based Handshake DNS resolver with JSON-RPC API
 * Optimized for Version 1 (V1) with lightweight sync & DNS resolution
 * Added JSON-RPC API for resolver integration with spvproc.go
 * Copyright (c) 2017-2018, Christopher Jeffrey (MIT License)
 * Modified for V1 optimization
 */

'use strict';

const assert = require('bsert');
const Chain = require('../blockchain/chain');
const Pool = require('../net/pool');
const Node = require('./node');
const HTTP = require('./http');
const RPC = require('./rpc');
const pkg = require('../pkg');
const { RootServer, RecursiveServer } = require('../dns/server');
const express = require('express');

/**
 * SPV Node - Optimized for lightweight DNS resolution with JSON-RPC API
 * @extends Node
 */
class SPVNode extends Node {
  /**
   * Create SPV node.
   * @constructor
   * @param {Object?} options
   */
  constructor(options) {
    super(pkg.core, pkg.cfg, 'debug.log');
    this.options = options || {};
    this.chain = new Chain({ prune: true, memory: false, spv: true }); // Lightweight sync
    this.pool = new Pool(this.options);
    this.http = new HTTP(this);
    this.rpc = new RPC(this);
    this.dns = new RootServer(this); // Handshake DNS server
    this.app = express();
    this.setupAPI();
  }

  /**
   * Open the SPV node.
   */
  async open() {
    await this.chain.open();
    await this.pool.open();
    await this.http.open();
    await this.rpc.open();
    await this.dns.open();
  }

  /**
   * Close the SPV node.
   */
  async close() {
    await this.dns.close();
    await this.rpc.close();
    await this.http.close();
    await this.pool.close();
    await this.chain.close();
  }

  /**
   * Setup JSON-RPC API for DNS resolution
   */
  setupAPI() {
    this.app.use(express.json());

    this.app.post('/resolve', async (req, res) => {
      const { name } = req.body;
      if (!name) return res.status(400).json({ error: 'Missing domain name' });
      try {
        const result = await this.dns.resolve(name);
        res.json({ name, result });
      } catch (err) {
        res.status(500).json({ error: err.message });
      }
    });

    this.app.get('/status', (req, res) => {
      res.json({ synced: this.chain.synced, height: this.chain.height });
    });

    this.app.post('/sync', async (req, res) => {
      try {
        await this.chain.sync();
        res.json({ message: 'Sync started' });
      } catch (err) {
        res.status(500).json({ error: err.message });
      }
    });

    this.app.listen(3000, () => console.log('SPV Resolver API running on port 3000'));
  }
}

module.exports = SPVNode;
