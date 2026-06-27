#pragma once

#include <string>
#include <vector>
#include <functional>
#include <iostream>
#include <chrono>
#include <sstream>

namespace cnext {
namespace test {

struct TestCase {
    std::string name;
    std::string suite;
    std::function<void()> func;
};

class TestRegistry {
public:
    static TestRegistry& Instance() {
        static TestRegistry instance;
        return instance;
    }
    
    void Register(const std::string& suite, const std::string& name, std::function<void()> func) {
        cases.push_back({name, suite, func});
    }
    
    const std::vector<TestCase>& Cases() const {
        return cases;
    }
    
    void Clear() {
        cases.clear();
    }
    
private:
    std::vector<TestCase> cases;
};

class TestRunner {
public:
    int Run(const std::string& filter = "") {
        int passed = 0;
        int failed = 0;
        int skipped = 0;
        double totalTime = 0.0;
        
        auto& registry = TestRegistry::Instance();
        auto& cases = registry.Cases();
        
        std::cout << "Running " << cases.size() << " tests..." << std::endl;
        std::cout << std::endl;
        
        for (const auto& tc : cases) {
            if (!filter.empty() && 
                tc.suite.find(filter) == std::string::npos && 
                tc.name.find(filter) == std::string::npos) {
                continue;
            }
            
            auto start = std::chrono::high_resolution_clock::now();
            bool testFailed = false;
            std::string errorMsg;
            
            try {
                tc.func();
            } catch (const std::exception& e) {
                testFailed = true;
                errorMsg = e.what();
            } catch (...) {
                testFailed = true;
                errorMsg = "unknown exception";
            }
            
            auto end = std::chrono::high_resolution_clock::now();
            auto duration = std::chrono::duration_cast<std::chrono::microseconds>(end - start);
            double ms = duration.count() / 1000.0;
            totalTime += ms;
            
            if (testFailed) {
                std::cout << "  FAIL: " << tc.suite << "::" << tc.name 
                         << " (" << ms << "ms)" << std::endl;
                std::cout << "    Error: " << errorMsg << std::endl;
                failed++;
            } else {
                std::cout << "  PASS: " << tc.suite << "::" << tc.name 
                         << " (" << ms << "ms)" << std::endl;
                passed++;
            }
        }
        
        std::cout << std::endl;
        std::cout << passed << " passed, " << failed << " failed, " 
                 << skipped << " skipped" << std::endl;
        std::cout << "Total time: " << totalTime << "ms" << std::endl;
        
        return failed > 0 ? 1 : 0;
    }
};

} // namespace test
} // namespace cnext

#define TEST(suite, name) \
    void suite##_##name(); \
    namespace { \
        struct Suite##_##name##Registrar { \
            Suite##_##name##Registrar() { \
                ::cnext::test::TestRegistry::Instance().Register(#suite, #name, suite##_##name); \
            } \
        } suite##_##name##registrar; \
    } \
    void suite##_##name()

#define TEST_F(suite, name) TEST(suite, name)

#define EXPECT_EQ(a, b) \
    do { \
        if ((a) != (b)) { \
            std::ostringstream oss; \
            oss << "Expected " #a " == " #b " (" << (a) << " != " << (b) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define EXPECT_NE(a, b) \
    do { \
        if ((a) == (b)) { \
            std::ostringstream oss; \
            oss << "Expected " #a " != " #b " (" << (a) << " == " << (b) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define EXPECT_TRUE(cond) \
    do { \
        if (!(cond)) { \
            throw std::runtime_error("Expected " #cond " to be true"); \
        } \
    } while(0)

#define EXPECT_FALSE(cond) \
    do { \
        if (cond) { \
            throw std::runtime_error("Expected " #cond " to be false"); \
        } \
    } while(0)

#define EXPECT_GT(a, b) \
    do { \
        if (!((a) > (b))) { \
            std::ostringstream oss; \
            oss << "Expected " #a " > " #b " (" << (a) << " <= " << (b) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define EXPECT_LT(a, b) \
    do { \
        if (!((a) < (b))) { \
            std::ostringstream oss; \
            oss << "Expected " #a " < " #b " (" << (a) << " >= " << (b) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define EXPECT_GE(a, b) \
    do { \
        if (!((a) >= (b))) { \
            std::ostringstream oss; \
            oss << "Expected " #a " >= " #b " (" << (a) << " < " << (b) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define EXPECT_LE(a, b) \
    do { \
        if (!((a) <= (b))) { \
            std::ostringstream oss; \
            oss << "Expected " #a " <= " #b " (" << (a) << " > " << (b) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define EXPECT_NEAR(a, b, epsilon) \
    do { \
        if (std::abs((a) - (b)) > (epsilon)) { \
            std::ostringstream oss; \
            oss << "Expected " #a " ~= " #b " (epsilon=" #epsilon ", actual=" << std::abs((a) - (b)) << ")"; \
            throw std::runtime_error(oss.str()); \
        } \
    } while(0)

#define ASSERT_EQ(a, b) EXPECT_EQ(a, b)
#define ASSERT_NE(a, b) EXPECT_NE(a, b)
#define ASSERT_TRUE(cond) EXPECT_TRUE(cond)
#define ASSERT_FALSE(cond) EXPECT_FALSE(cond)
#define ASSERT_GT(a, b) EXPECT_GT(a, b)
#define ASSERT_LT(a, b) EXPECT_LT(a, b)
#define ASSERT_GE(a, b) EXPECT_GE(a, b)
#define ASSERT_LE(a, b) EXPECT_LE(a, b)
