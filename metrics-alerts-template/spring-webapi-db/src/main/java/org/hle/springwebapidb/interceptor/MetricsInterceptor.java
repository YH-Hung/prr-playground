package org.hle.springwebapidb.interceptor;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

@Component
public class MetricsInterceptor implements HandlerInterceptor {
    
    private static final Logger logger = LoggerFactory.getLogger(MetricsInterceptor.class);
    private final MeterRegistry meterRegistry;
    private static final String TIMER_NAME = "http.server.requests";
    private static final String COUNTER_NAME = "http.server.requests.total";
    
    public MetricsInterceptor(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
    }
    
    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        request.setAttribute("startTime", System.currentTimeMillis());
        return true;
    }
    
    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) {
        Long startTime = (Long) request.getAttribute("startTime");
        if (startTime != null) {
            long duration = System.currentTimeMillis() - startTime;
            String method = request.getMethod();
            String uri = request.getRequestURI();
            int status = response.getStatus();
            
            // Record timer with tags
            Timer.Sample sample = Timer.start(meterRegistry);
            sample.stop(Timer.builder(TIMER_NAME)
                    .tag("method", method)
                    .tag("uri", sanitizeUri(uri))
                    .tag("status", String.valueOf(status))
                    .register(meterRegistry));
            
            // Record counter for total requests
            Counter.builder(COUNTER_NAME)
                    .tag("method", method)
                    .tag("uri", sanitizeUri(uri))
                    .tag("status", String.valueOf(status))
                    .register(meterRegistry)
                    .increment();
            
            // Record error counter if status is 5xx
            if (status >= 500) {
                Counter.builder("http.server.errors")
                        .tag("method", method)
                        .tag("uri", sanitizeUri(uri))
                        .tag("status", String.valueOf(status))
                        .register(meterRegistry)
                        .increment();
                logger.warn("HTTP error: {} {} - Status: {} - Duration: {}ms", method, uri, status, duration);
            } else if (status >= 400) {
                Counter.builder("http.server.client.errors")
                        .tag("method", method)
                        .tag("uri", sanitizeUri(uri))
                        .tag("status", String.valueOf(status))
                        .register(meterRegistry)
                        .increment();
            }
            
            // Log slow requests
            if (duration > 1000) {
                logger.warn("Slow request detected: {} {} - Duration: {}ms", method, uri, duration);
            }
        }
    }
    
    private String sanitizeUri(String uri) {
        // Replace path variables with placeholders for better metric aggregation
        if (uri.startsWith("/api/users/") && uri.matches("/api/users/\\d+")) {
            return "/api/users/{id}";
        }
        if (uri.startsWith("/api/users/email/")) {
            return "/api/users/email/{email}";
        }
        if (uri.startsWith("/api/users/status/")) {
            return "/api/users/status/{status}";
        }
        if (uri.startsWith("/api/users/external/")) {
            return "/api/users/external/{serviceName}";
        }
        return uri;
    }
}

