add_executable(vm_tests all_expressions_test.cpp)

target_link_libraries(vm_tests smalltalk_vm)

add_test(NAME vm_tests COMMAND vm_tests)