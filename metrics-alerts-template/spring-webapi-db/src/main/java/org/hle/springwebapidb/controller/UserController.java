package org.hle.springwebapidb.controller;

import org.hle.springwebapidb.entity.User;
import org.hle.springwebapidb.service.UserService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/users")
public class UserController {
    
    private static final Logger logger = LoggerFactory.getLogger(UserController.class);
    private final UserService userService;
    
    public UserController(UserService userService) {
        this.userService = userService;
    }
    
    @PostMapping
    public ResponseEntity<User> createUser(@RequestBody User user) {
        logger.info("Creating user with email: {}", user.getEmail());
        try {
            User createdUser = userService.createUser(user);
            return ResponseEntity.status(HttpStatus.CREATED).body(createdUser);
        } catch (Exception e) {
            logger.error("Error creating user", e);
            throw e;
        }
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<User> getUserById(@PathVariable Long id) {
        logger.info("Getting user with id: {}", id);
        return userService.getUserById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping
    public ResponseEntity<List<User>> getAllUsers() {
        logger.info("Getting all users");
        return ResponseEntity.ok(userService.getAllUsers());
    }
    
    @GetMapping("/email/{email}")
    public ResponseEntity<User> getUserByEmail(@PathVariable String email) {
        logger.info("Getting user with email: {}", email);
        return userService.getUserByEmail(email)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<User> updateUser(@PathVariable Long id, @RequestBody User userDetails) {
        logger.info("Updating user with id: {}", id);
        try {
            User updatedUser = userService.updateUser(id, userDetails);
            return ResponseEntity.ok(updatedUser);
        } catch (Exception e) {
            logger.error("Error updating user", e);
            throw e;
        }
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteUser(@PathVariable Long id) {
        logger.info("Deleting user with id: {}", id);
        try {
            userService.deleteUser(id);
            return ResponseEntity.noContent().build();
        } catch (Exception e) {
            logger.error("Error deleting user", e);
            throw e;
        }
    }
    
    @GetMapping("/status/{status}")
    public ResponseEntity<List<User>> getUsersByStatus(@PathVariable String status) {
        logger.info("Getting users with status: {}", status);
        return ResponseEntity.ok(userService.getUsersByStatus(status));
    }
    
    @GetMapping("/status/{status}/count")
    public ResponseEntity<Map<String, Object>> countUsersByStatus(@PathVariable String status) {
        logger.info("Counting users with status: {}", status);
        long count = userService.countUsersByStatus(status);
        return ResponseEntity.ok(Map.of("status", status, "count", count));
    }
    
    // Endpoint to simulate external service call
    @GetMapping("/external/{serviceName}")
    public ResponseEntity<Map<String, String>> callExternalService(@PathVariable String serviceName) {
        logger.info("Calling external service: {}", serviceName);
        try {
            String result = userService.callExternalService(serviceName);
            return ResponseEntity.ok(Map.of("service", serviceName, "result", result));
        } catch (Exception e) {
            logger.error("Error calling external service", e);
            throw e;
        }
    }
    
    // Endpoint to simulate failure scenarios for testing alerts
    @GetMapping("/test/error")
    public ResponseEntity<Map<String, String>> triggerError() {
        logger.warn("Triggering test error endpoint");
        throw new RuntimeException("Test error for alerting demonstration");
    }
    
    @GetMapping("/test/slow")
    public ResponseEntity<Map<String, String>> triggerSlowResponse() throws InterruptedException {
        logger.warn("Triggering slow response endpoint");
        Thread.sleep(2000); // 2 second delay
        return ResponseEntity.ok(Map.of("message", "Slow response completed"));
    }
}

