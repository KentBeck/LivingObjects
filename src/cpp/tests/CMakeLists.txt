# Fetch and build Google Test
include(FetchContent)
FetchContent_Declare(
  googletest
  GIT_REPOSITORY https://github.com/google/googletest.git
  GIT_TAG release-1.12.1
)
# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

# Tests
file(GLOB_RECURSE TEST_SOURCES "*.cpp")
add_executable(vm_tests ${TEST_SOURCES})
target_link_libraries(vm_tests
  smalltalk_vm
  gtest
  gtest_main
)

# Register tests
include(GoogleTest)
gtest_discover_tests(vm_tests)