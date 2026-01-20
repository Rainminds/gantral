import unittest
import os
import sys

# Add parent dir to path to find runner module
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from secret_store.resolver import SecretResolver

class TestSecretResolver(unittest.TestCase):
    def setUp(self):
        self.resolver = SecretResolver()
        os.environ["TEST_ENV_SECRET"] = "secret-value"

    def test_pass_through(self):
        """Should pass non-secret strings through unchanged."""
        self.assertEqual(self.resolver.resolve("plain-text"), "plain-text")
        self.assertEqual(self.resolver.resolve(123), 123)

    def test_resolve_env(self):
        """Should resolve gantral+secret://env/..."""
        uri = "gantral+secret://env/TEST_ENV_SECRET"
        self.assertEqual(self.resolver.resolve(uri), "secret-value")

    def test_resolve_vault_mock(self):
        """Should resolve gantral+secret://vault/..."""
        uri = "gantral+secret://vault/my-secret"
        # MockVaultProvider returns "mock-vault-value-for-{path}"
        expected = "mock-vault-value-for-my-secret"
        self.assertEqual(self.resolver.resolve(uri), expected)
    
    def test_unknown_provider(self):
        """Should return empty string for unknown provider."""
        uri = "gantral+secret://unknown/path"
        self.assertEqual(self.resolver.resolve(uri), "")

if __name__ == '__main__':
    unittest.main()
