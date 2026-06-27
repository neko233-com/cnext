#pragma once

#include <cmath>
#include <stdexcept>
#include <string>
#include <sstream>

namespace cnext {
namespace test {

template<typename T>
void AssertEqual(const T& a, const T& b, const std::string& msg = "") {
    if (a != b) {
        std::ostringstream oss;
        oss << "AssertEqual failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

template<typename T>
void AssertNotEqual(const T& a, const T& b, const std::string& msg = "") {
    if (a == b) {
        std::ostringstream oss;
        oss << "AssertNotEqual failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

template<typename T>
void AssertNear(const T& a, const T& b, const T& epsilon, const std::string& msg = "") {
    if (std::abs(a - b) > epsilon) {
        std::ostringstream oss;
        oss << "AssertNear failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

inline void AssertTrue(bool condition, const std::string& msg = "") {
    if (!condition) {
        std::ostringstream oss;
        oss << "AssertTrue failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

inline void AssertFalse(bool condition, const std::string& msg = "") {
    if (condition) {
        std::ostringstream oss;
        oss << "AssertFalse failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

template<typename T>
void AssertGreaterThan(const T& a, const T& b, const std::string& msg = "") {
    if (!(a > b)) {
        std::ostringstream oss;
        oss << "AssertGreaterThan failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

template<typename T>
void AssertLessThan(const T& a, const T& b, const std::string& msg = "") {
    if (!(a < b)) {
        std::ostringstream oss;
        oss << "AssertLessThan failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

template<typename T>
void AssertGreaterOrEqual(const T& a, const T& b, const std::string& msg = "") {
    if (!(a >= b)) {
        std::ostringstream oss;
        oss << "AssertGreaterOrEqual failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

template<typename T>
void AssertLessOrEqual(const T& a, const T& b, const std::string& msg = "") {
    if (!(a <= b)) {
        std::ostringstream oss;
        oss << "AssertLessOrEqual failed";
        if (!msg.empty()) {
            oss << ": " << msg;
        }
        throw std::runtime_error(oss.str());
    }
}

} // namespace test
} // namespace cnext
