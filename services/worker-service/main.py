#!/usr/bin/env python3
# services/worker-service/main.py
# This simulates telex_queue_processor - processes background jobs

import time
import random
import logging
from datetime import datetime
from flask import Flask, jsonify, request

# Basic logging setup (we'll improve this later)
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] [%(levelname)s] %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)

app = Flask(__name__)
logger = logging.getLogger(__name__)

# Simulate job queue
job_queue = []
processed_jobs = []


@app.route('/health')
def health():
    """Health check endpoint"""
    logger.info("Health check requested")
    return jsonify({
        "status": "healthy",
        "service": "worker-service",
        "queue_size": len(job_queue),
        "processed": len(processed_jobs)
    })


@app.route('/jobs', methods=['GET', 'POST'])
def jobs():
    """Job management endpoint"""
    if request.method == 'GET':
        logger.info(f"Fetching jobs - Queue: {len(job_queue)}, Processed: {len(processed_jobs)}")
        return jsonify({
            "queued": job_queue,
            "processed": processed_jobs
        })
    
    elif request.method == 'POST':
        data = request.get_json()
        job = {
            "id": f"job_{random.randint(1000, 9999)}",
            "type": data.get("type", "process_message"),
            "data": data.get("data", {}),
            "created_at": datetime.now().isoformat(),
            "status": "queued"
        }
        job_queue.append(job)
        logger.info(f"Job queued: {job['id']} - Type: {job['type']}")
        return jsonify(job), 201


@app.route('/jobs/process', methods=['POST'])
def process_jobs():
    """Process jobs in queue"""
    if not job_queue:
        logger.warning("Process triggered but queue is empty")
        return jsonify({"message": "Queue is empty"}), 200
    
    job = job_queue.pop(0)
    logger.info(f"Processing job: {job['id']}")
    
    # Simulate processing
    start_time = time.time()
    processing_time = random.uniform(0.1, 2.0)
    time.sleep(processing_time)
    
    # Simulate random failures (15% chance)
    if random.random() < 0.15:
        logger.error(f"Job {job['id']} failed - Processing error")
        job['status'] = 'failed'
        job['error'] = 'Random processing error'
    else:
        logger.info(f"Job {job['id']} completed successfully in {processing_time:.2f}s")
        job['status'] = 'completed'
    
    job['processed_at'] = datetime.now().isoformat()
    job['duration'] = time.time() - start_time
    processed_jobs.append(job)
    
    return jsonify(job)


@app.route('/jobs/stats')
def job_stats():
    """Get job statistics"""
    completed = sum(1 for j in processed_jobs if j['status'] == 'completed')
    failed = sum(1 for j in processed_jobs if j['status'] == 'failed')
    
    logger.info(f"Stats requested - Completed: {completed}, Failed: {failed}")
    
    return jsonify({
        "queue_size": len(job_queue),
        "processed_total": len(processed_jobs),
        "completed": completed,
        "failed": failed,
        "success_rate": completed / len(processed_jobs) if processed_jobs else 0
    })


@app.route('/chaos')
def chaos():
    """Chaos endpoint for testing"""
    action = request.args.get('action', 'status')
    
    if action == 'slow':
        logger.warning("Chaos mode: Slow processing triggered")
        time.sleep(10)
        return jsonify({"message": "Slow processing completed"})
    
    elif action == 'error':
        logger.error("Chaos mode: Error triggered")
        return jsonify({"error": "Chaos error"}), 500
    
    elif action == 'memory':
        logger.warning("Chaos mode: Memory leak simulation")
        # Don't actually do this in production!
        big_list = [random.random() for _ in range(1000000)]
        return jsonify({"message": f"Allocated {len(big_list)} items"})
    
    else:
        return jsonify({"message": "Chaos mode available"})


if __name__ == '__main__':
    logger.info("Starting Worker Service on :8001")
    app.run(host='0.0.0.0', port=8001, debug=False)