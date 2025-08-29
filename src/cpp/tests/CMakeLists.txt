add_executable(vm_tests all_expressions_test.cpp)
target_link_libraries(vm_tests smalltalk_vm)
add_test(NAME vm_tests COMMAND vm_tests)

# Image build and boot test: builds a binary image from .st sources,
# loads it, and runs a simple expression to ensure VM + primitives work.
add_executable(image_build_and_boot_test image_build_and_boot_test.cpp)
target_link_libraries(image_build_and_boot_test smalltalk_vm)
add_test(NAME image_build_and_boot_test COMMAND image_build_and_boot_test)
