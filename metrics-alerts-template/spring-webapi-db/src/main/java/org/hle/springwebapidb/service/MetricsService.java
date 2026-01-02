package org.hle.springwebapidb.service;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import jakarta.annotation.PostConstruct;
import org.springframework.stereotype.Service;

import java.util.concurrent.atomic.AtomicInteger;

@Service
public class MetricsService {
    
    private final MeterRegistry meterRegistry;
    private final AtomicInteger activeOperations = new AtomicInteger(0);
    private Counter userCreatedCounter;
    private Counter userUpdatedCounter;
    private Counter userDeletedCounter;
    private Timer userOperationTimer;
    
    public MetricsService(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
    }
    
    @PostConstruct
    public void init() {
        // Counters for business events
        userCreatedCounter = Counter.builder("custom.user.created")
                .description("Total number of users created")
                .tag("operation", "create")
                .register(meterRegistry);
        
        userUpdatedCounter = Counter.builder("custom.user.updated")
                .description("Total number of users updated")
                .tag("operation", "update")
                .register(meterRegistry);
        
        userDeletedCounter = Counter.builder("custom.user.deleted")
                .description("Total number of users deleted")
                .tag("operation", "delete")
                .register(meterRegistry);
        
        // Timer for operation duration
        userOperationTimer = Timer.builder("custom.user.operation.duration")
                .description("Duration of user operations")
                .register(meterRegistry);
        
        // Gauge for active operations
        Gauge.builder("custom.user.active.operations", activeOperations, AtomicInteger::get)
                .description("Number of active user operations")
                .register(meterRegistry);
    }
    
    public void incrementUserCreated() {
        userCreatedCounter.increment();
    }
    
    public void incrementUserUpdated() {
        userUpdatedCounter.increment();
    }
    
    public void incrementUserDeleted() {
        userDeletedCounter.increment();
    }
    
    public void incrementUserOperationErrors(String errorType) {
        Counter.builder("custom.user.operation.errors")
                .tag("error.type", errorType)
                .register(meterRegistry)
                .increment();
    }
    
    public Timer.Sample startUserOperationTimer() {
        activeOperations.incrementAndGet();
        return Timer.start(meterRegistry);
    }
    
    public void stopUserOperationTimer(Timer.Sample sample) {
        sample.stop(userOperationTimer);
        activeOperations.decrementAndGet();
    }
    
    public void recordExternalCallDuration(String service, long durationMs) {
        Timer.builder("custom.external.call.duration")
                .tag("service", service)
                .register(meterRegistry)
                .record(java.time.Duration.ofMillis(durationMs));
    }
    
    public void incrementExternalCallErrors(String service) {
        Counter.builder("custom.external.call.errors")
                .tag("service", service)
                .register(meterRegistry)
                .increment();
    }
}

