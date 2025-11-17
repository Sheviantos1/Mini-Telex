// services/webhook-service/index.js
// This simulates n8n-node-executor - handles webhooks and workflows

const express = require('express');
const app = express();
const port = 8002;

app.use(express.json());

// Simple in-memory storage
const webhooks = [];
const executions = [];

// Basic logging (we'll improve this later)
const log = (level, message, data = {}) => {
  const timestamp = new Date().toISOString();
  console.log(`[${timestamp}] [${level}] ${message}`, data);
};

// Health check
app.get('/health', (req, res) => {
  log('INFO', 'Health check requested');
  res.json({
    status: 'healthy',
    service: 'webhook-service',
    webhooks: webhooks.length,
    executions: executions.length
  });
});

// List webhooks
app.get('/webhooks', (req, res) => {
  log('INFO', `Fetching webhooks - Total: ${webhooks.length}`);
  res.json(webhooks);
});

// Create webhook
app.post('/webhooks', (req, res) => {
  const webhook = {
    id: `webhook_${Math.floor(Math.random() * 10000)}`,
    name: req.body.name || 'Unnamed Webhook',
    url: `/webhooks/${Math.random().toString(36).substring(7)}`,
    workflow: req.body.workflow || 'default',
    created_at: new Date().toISOString(),
    active: true
  };
  
  webhooks.push(webhook);
  log('INFO', `Webhook created: ${webhook.id}`, { name: webhook.name });
  
  res.status(201).json(webhook);
});

// Trigger webhook
app.post('/webhooks/:id/trigger', (req, res) => {
  const webhookId = req.params.id;
  const webhook = webhooks.find(w => w.id === webhookId);
  
  if (!webhook) {
    log('ERROR', `Webhook not found: ${webhookId}`);
    return res.status(404).json({ error: 'Webhook not found' });
  }
  
  log('INFO', `Webhook triggered: ${webhookId}`, { workflow: webhook.workflow });
  
  // Simulate workflow execution
  const execution = {
    id: `exec_${Math.floor(Math.random() * 10000)}`,
    webhook_id: webhookId,
    workflow: webhook.workflow,
    started_at: new Date().toISOString(),
    status: 'running'
  };
  
  executions.push(execution);
  
  // Simulate processing
  const processingTime = Math.random() * 3000; // 0-3 seconds
  
  setTimeout(() => {
    // Random failure (10% chance)
    if (Math.random() < 0.1) {
      execution.status = 'failed';
      execution.error = 'Random workflow error';
      log('ERROR', `Workflow execution failed: ${execution.id}`, { error: execution.error });
    } else {
      execution.status = 'completed';
      log('INFO', `Workflow execution completed: ${execution.id}`, { duration: processingTime });
    }
    
    execution.completed_at = new Date().toISOString();
    execution.duration_ms = processingTime;
  }, processingTime);
  
  res.json({
    message: 'Webhook triggered',
    execution_id: execution.id
  });
});

// List executions
app.get('/executions', (req, res) => {
  log('INFO', `Fetching executions - Total: ${executions.length}`);
  res.json(executions);
});

// Get execution details
app.get('/executions/:id', (req, res) => {
  const execution = executions.find(e => e.id === req.params.id);
  
  if (!execution) {
    log('ERROR', `Execution not found: ${req.params.id}`);
    return res.status(404).json({ error: 'Execution not found' });
  }
  
  log('INFO', `Execution details requested: ${execution.id}`);
  res.json(execution);
});

// Execution statistics
app.get('/executions/stats', (req, res) => {
  const completed = executions.filter(e => e.status === 'completed').length;
  const failed = executions.filter(e => e.status === 'failed').length;
  const running = executions.filter(e => e.status === 'running').length;
  
  log('INFO', 'Execution stats requested', { completed, failed, running });
  
  res.json({
    total: executions.length,
    completed,
    failed,
    running,
    success_rate: executions.length > 0 ? (completed / executions.length) : 0
  });
});

// Chaos endpoint for testing
app.get('/chaos', (req, res) => {
  const action = req.query.action || 'status';
  
  switch (action) {
    case 'slow':
      log('WARN', 'Chaos mode: Slow response triggered');
      setTimeout(() => {
        res.json({ message: 'Slow response completed' });
      }, 5000);
      break;
      
    case 'error':
      log('ERROR', 'Chaos mode: Error triggered');
      res.status(500).json({ error: 'Chaos error' });
      break;
      
    case 'crash':
      log('FATAL', 'Chaos mode: Crash triggered');
      throw new Error('Chaos crash!');
      
    default:
      res.json({ message: 'Chaos mode available. Use ?action=slow|error|crash' });
  }
});

// Error handling
app.use((err, req, res, next) => {
  log('ERROR', 'Unhandled error', { error: err.message, stack: err.stack });
  res.status(500).json({ error: 'Internal server error' });
});

// Start server
app.listen(port, () => {
  log('INFO', `Webhook Service started on port ${port}`);
});