"""
Executable test suite for integer primitive methods.
These tests define the expected behavior that all VM implementations must match.
"""

import pytest
from typing import Any, List


class VMInterface:
    """
    Interface that VM implementations must provide for testing.
    Adapt your VM to this interface to run the test suite.
    """

    def push(self, value: Any) -> None:
        """Push a value onto the stack"""
        raise NotImplementedError

    def pop(self) -> Any:
        """Pop and return top value from stack"""
        raise NotImplementedError

    def call_primitive(self, primitive_number: int) -> Any:
        """
        Call a primitive method.
        Receiver and arguments should already be on stack.
        Returns the result (also pushed onto stack).
        """
        raise NotImplementedError

    @property
    def stack(self) -> List[Any]:
        """Return current stack contents (bottom to top)"""
        raise NotImplementedError

    def reset(self) -> None:
        """Reset VM state for next test"""
        raise NotImplementedError


class TestPrimitive1_Add:
    """Tests for Primitive 1: SmallInteger>>+"""

    def test_basic_addition(self, vm: VMInterface):
        """3 + 4 should return 7"""
        vm.push(3)  # receiver
        vm.push(4)  # argument
        result = vm.call_primitive(1)

        assert result == 7
        assert vm.stack == [7]

    def test_zero_addition(self, vm: VMInterface):
        """0 + 0 should return 0"""
        vm.push(0)
        vm.push(0)
        result = vm.call_primitive(1)

        assert result == 0
        assert vm.stack == [0]

    def test_negative_addition(self, vm: VMInterface):
        """-5 + 3 should return -2"""
        vm.push(-5)
        vm.push(3)
        result = vm.call_primitive(1)

        assert result == -2
        assert vm.stack == [-2]

    def test_large_numbers(self, vm: VMInterface):
        """1000000 + 2000000 should return 3000000"""
        vm.push(1000000)
        vm.push(2000000)
        result = vm.call_primitive(1)

        assert result == 3000000
        assert vm.stack == [3000000]

    def test_commutative(self, vm: VMInterface):
        """Addition should be commutative: a + b == b + a"""
        # Test 3 + 4
        vm.push(3)
        vm.push(4)
        result1 = vm.call_primitive(1)
        vm.reset()

        # Test 4 + 3
        vm.push(4)
        vm.push(3)
        result2 = vm.call_primitive(1)

        assert result1 == result2 == 7

    def test_identity_element(self, vm: VMInterface):
        """0 is the identity element: n + 0 == n"""
        vm.push(42)
        vm.push(0)
        result = vm.call_primitive(1)

        assert result == 42

    def test_type_error_string_receiver(self, vm: VMInterface):
        """Should fail when receiver is not SmallInteger"""
        vm.push("string")
        vm.push(4)

        with pytest.raises(TypeError):
            vm.call_primitive(1)

    def test_type_error_string_argument(self, vm: VMInterface):
        """Should fail when argument is not SmallInteger"""
        vm.push(3)
        vm.push("string")

        with pytest.raises(TypeError):
            vm.call_primitive(1)

    def test_overflow_positive(self, vm: VMInterface):
        """Should fail on positive overflow (31-bit limit)"""
        max_int = 2**30 - 1  # Maximum 31-bit signed integer
        vm.push(max_int)
        vm.push(1)

        with pytest.raises(OverflowError):
            vm.call_primitive(1)

    def test_overflow_negative(self, vm: VMInterface):
        """Should fail on negative overflow (31-bit limit)"""
        min_int = -(2**30)  # Minimum 31-bit signed integer
        vm.push(min_int)
        vm.push(-1)

        with pytest.raises(OverflowError):
            vm.call_primitive(1)

    def test_stack_effect(self, vm: VMInterface):
        """Should pop 2 values and push 1 (net effect: -1)"""
        vm.push(1)
        vm.push(2)
        vm.push(3)  # receiver
        vm.push(4)  # argument

        initial_depth = len(vm.stack)
        vm.call_primitive(1)
        final_depth = len(vm.stack)

        assert final_depth == initial_depth - 1
        assert vm.stack == [1, 2, 7]


class TestPrimitive2_Subtract:
    """Tests for Primitive 2: SmallInteger>>-"""

    def test_basic_subtraction(self, vm: VMInterface):
        """10 - 3 should return 7"""
        vm.push(10)
        vm.push(3)
        result = vm.call_primitive(2)

        assert result == 7
        assert vm.stack == [7]

    def test_zero_result(self, vm: VMInterface):
        """5 - 5 should return 0"""
        vm.push(5)
        vm.push(5)
        result = vm.call_primitive(2)

        assert result == 0

    def test_negative_result(self, vm: VMInterface):
        """3 - 10 should return -7"""
        vm.push(3)
        vm.push(10)
        result = vm.call_primitive(2)

        assert result == -7

    def test_identity_element(self, vm: VMInterface):
        """0 is the identity element: n - 0 == n"""
        vm.push(42)
        vm.push(0)
        result = vm.call_primitive(2)

        assert result == 42

    def test_inverse_of_addition(self, vm: VMInterface):
        """Subtraction is inverse of addition: (a + b) - b == a"""
        # First: 3 + 4 = 7
        vm.push(3)
        vm.push(4)
        sum_result = vm.call_primitive(1)
        vm.reset()

        # Then: 7 - 4 = 3
        vm.push(sum_result)
        vm.push(4)
        diff_result = vm.call_primitive(2)

        assert diff_result == 3


class TestPrimitive3_LessThan:
    """Tests for Primitive 3: SmallInteger>><"""

    def test_less_than_true(self, vm: VMInterface):
        """3 < 5 should return true"""
        vm.push(3)
        vm.push(5)
        result = vm.call_primitive(3)

        assert result is True
        assert vm.stack == [True]

    def test_less_than_false(self, vm: VMInterface):
        """7 < 2 should return false"""
        vm.push(7)
        vm.push(2)
        result = vm.call_primitive(3)

        assert result is False

    def test_equal_values(self, vm: VMInterface):
        """5 < 5 should return false"""
        vm.push(5)
        vm.push(5)
        result = vm.call_primitive(3)

        assert result is False

    def test_negative_numbers(self, vm: VMInterface):
        """-10 < -5 should return true"""
        vm.push(-10)
        vm.push(-5)
        result = vm.call_primitive(3)

        assert result is True

    def test_transitive(self, vm: VMInterface):
        """If a < b and b < c, then a < c"""
        # 1 < 2
        vm.push(1)
        vm.push(2)
        result1 = vm.call_primitive(3)
        vm.reset()

        # 2 < 3
        vm.push(2)
        vm.push(3)
        result2 = vm.call_primitive(3)
        vm.reset()

        # 1 < 3
        vm.push(1)
        vm.push(3)
        result3 = vm.call_primitive(3)

        assert result1 and result2 and result3


class TestPrimitive10_Divide:
    """Tests for Primitive 10: SmallInteger>>/"""

    def test_exact_division(self, vm: VMInterface):
        """20 / 4 should return 5"""
        vm.push(20)
        vm.push(4)
        result = vm.call_primitive(10)

        assert result == 5

    def test_truncated_division(self, vm: VMInterface):
        """7 / 2 should return 3 (truncated toward zero)"""
        vm.push(7)
        vm.push(2)
        result = vm.call_primitive(10)

        assert result == 3

    def test_negative_truncated_division(self, vm: VMInterface):
        """-7 / 2 should return -3 (truncated toward zero)"""
        vm.push(-7)
        vm.push(2)
        result = vm.call_primitive(10)

        assert result == -3

    def test_division_by_one(self, vm: VMInterface):
        """n / 1 should return n"""
        vm.push(42)
        vm.push(1)
        result = vm.call_primitive(10)

        assert result == 42

    def test_division_by_zero(self, vm: VMInterface):
        """10 / 0 should raise ZeroDivisionError"""
        vm.push(10)
        vm.push(0)

        with pytest.raises(ZeroDivisionError):
            vm.call_primitive(10)

    def test_zero_divided_by_nonzero(self, vm: VMInterface):
        """0 / 5 should return 0"""
        vm.push(0)
        vm.push(5)
        result = vm.call_primitive(10)

        assert result == 0


# Pytest fixtures
@pytest.fixture
def vm():
    """
    Fixture that provides a VM instance for testing.
    Replace MockVM with your actual VM implementation.
    """
    vm_instance = MockVM()
    yield vm_instance
    vm_instance.reset()


class MockVM(VMInterface):
    """
    Mock VM implementation for demonstration.
    Replace this with your actual VM implementation.
    """

    def __init__(self):
        self._stack = []

    def push(self, value: Any) -> None:
        self._stack.append(value)

    def pop(self) -> Any:
        if not self._stack:
            raise RuntimeError("Stack underflow")
        return self._stack.pop()

    def call_primitive(self, primitive_number: int) -> Any:
        # This is a mock implementation
        # Your VM should implement actual primitive dispatch
        if primitive_number == 1:  # Add
            arg = self.pop()
            receiver = self.pop()
            if not isinstance(receiver, int) or not isinstance(arg, int):
                raise TypeError("Arguments must be integers")
            result = receiver + arg
            # Check overflow (31-bit signed)
            if result < -(2**30) or result >= 2**30:
                raise OverflowError("Integer overflow")
            self.push(result)
            return result
        # Add other primitives...
        raise NotImplementedError(f"Primitive {primitive_number} not implemented")

    @property
    def stack(self) -> List[Any]:
        return self._stack.copy()

    def reset(self) -> None:
        self._stack.clear()


if __name__ == "__main__":
    # Run tests with: python -m pytest test_integer_primitives.py -v
    pytest.main([__file__, "-v"])
