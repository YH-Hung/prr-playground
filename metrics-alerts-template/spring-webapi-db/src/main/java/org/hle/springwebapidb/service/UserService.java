package org.hle.springwebapidb.service;

import io.micrometer.core.instrument.Timer;
import org.hle.springwebapidb.entity.User;
import org.hle.springwebapidb.repository.UserRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;
import java.util.Random;

@Service
public class UserService {
    
    private final UserRepository userRepository;
    private final MetricsService metricsService;
    private final Random random = new Random();
    
    public UserService(UserRepository userRepository, MetricsService metricsService) {
        this.userRepository = userRepository;
        this.metricsService = metricsService;
    }
    
    @Transactional
    public User createUser(User user) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            if (userRepository.existsByEmail(user.getEmail())) {
                metricsService.incrementUserOperationErrors("duplicate_email");
                throw new IllegalArgumentException("User with email " + user.getEmail() + " already exists");
            }
            
            // Simulate occasional slow operations
            simulateRandomDelay();
            
            user.setStatus("ACTIVE");
            User savedUser = userRepository.save(user);
            metricsService.incrementUserCreated();
            return savedUser;
        } catch (Exception e) {
            metricsService.incrementUserOperationErrors("create_failed");
            throw e;
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional(readOnly = true)
    public Optional<User> getUserById(Long id) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            // Simulate occasional timeout scenarios
            if (random.nextInt(100) < 2) { // 2% chance
                Thread.sleep(3000); // Simulate timeout
            }
            
            return userRepository.findById(id);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            metricsService.incrementUserOperationErrors("timeout");
            throw new RuntimeException("Operation timed out", e);
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional(readOnly = true)
    public List<User> getAllUsers() {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            simulateRandomDelay();
            return userRepository.findAll();
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional(readOnly = true)
    public Optional<User> getUserByEmail(String email) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            return userRepository.findByEmail(email);
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional
    public User updateUser(Long id, User userDetails) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            User user = userRepository.findById(id)
                    .orElseThrow(() -> {
                        metricsService.incrementUserOperationErrors("not_found");
                        return new IllegalArgumentException("User not found with id: " + id);
                    });
            
            // Simulate occasional errors
            if (random.nextInt(100) < 1) { // 1% chance
                throw new RuntimeException("Simulated database error");
            }
            
            user.setName(userDetails.getName());
            user.setEmail(userDetails.getEmail());
            user.setStatus(userDetails.getStatus());
            
            User updatedUser = userRepository.save(user);
            metricsService.incrementUserUpdated();
            return updatedUser;
        } catch (Exception e) {
            metricsService.incrementUserOperationErrors("update_failed");
            throw e;
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional
    public void deleteUser(Long id) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            if (!userRepository.existsById(id)) {
                metricsService.incrementUserOperationErrors("not_found");
                throw new IllegalArgumentException("User not found with id: " + id);
            }
            
            userRepository.deleteById(id);
            metricsService.incrementUserDeleted();
        } catch (Exception e) {
            metricsService.incrementUserOperationErrors("delete_failed");
            throw e;
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional(readOnly = true)
    public List<User> getUsersByStatus(String status) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            return userRepository.findByStatus(status);
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    @Transactional(readOnly = true)
    public long countUsersByStatus(String status) {
        Timer.Sample sample = metricsService.startUserOperationTimer();
        try {
            return userRepository.countByStatus(status);
        } finally {
            metricsService.stopUserOperationTimer(sample);
        }
    }
    
    // Simulate external service call with metrics
    public String callExternalService(String serviceName) {
        long startTime = System.currentTimeMillis();
        try {
            // Simulate external call
            Thread.sleep(random.nextInt(500) + 100);
            
            // Simulate occasional failures
            if (random.nextInt(100) < 5) { // 5% failure rate
                metricsService.incrementExternalCallErrors(serviceName);
                throw new RuntimeException("External service " + serviceName + " failed");
            }
            
            long duration = System.currentTimeMillis() - startTime;
            metricsService.recordExternalCallDuration(serviceName, duration);
            return "Success from " + serviceName;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            metricsService.incrementExternalCallErrors(serviceName);
            throw new RuntimeException("External call interrupted", e);
        } catch (Exception e) {
            long duration = System.currentTimeMillis() - startTime;
            metricsService.recordExternalCallDuration(serviceName, duration);
            metricsService.incrementExternalCallErrors(serviceName);
            throw e;
        }
    }
    
    private void simulateRandomDelay() {
        try {
            // Random delay between 10-200ms
            Thread.sleep(random.nextInt(190) + 10);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}

