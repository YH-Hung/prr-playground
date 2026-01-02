package org.hle.springwebapidb.config;

import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.config.MeterFilter;
import jakarta.annotation.PostConstruct;
import org.springframework.context.annotation.Configuration;

@Configuration
public class MetricsConfig {
    
    private final MeterRegistry meterRegistry;
    
    public MetricsConfig(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
    }
    
    @PostConstruct
    public void configureMetrics() {
        meterRegistry.config()
                .commonTags("application", "spring-webapi-db")
                .meterFilter(MeterFilter.denyNameStartsWith("jvm.threads"))
                .meterFilter(MeterFilter.acceptNameStartsWith("http"))
                .meterFilter(MeterFilter.acceptNameStartsWith("jvm.memory"))
                .meterFilter(MeterFilter.acceptNameStartsWith("jvm.gc"))
                .meterFilter(MeterFilter.acceptNameStartsWith("process"))
                .meterFilter(MeterFilter.acceptNameStartsWith("system"))
                .meterFilter(MeterFilter.acceptNameStartsWith("hikari"))
                .meterFilter(MeterFilter.acceptNameStartsWith("jpa"))
                .meterFilter(MeterFilter.acceptNameStartsWith("custom"));
    }
}

