from abc import ABC, abstractmethod
import os
import logging

logger = logging.getLogger("gantral-runner")

class SecretProvider(ABC):
    @abstractmethod
    def get_secret(self, path: str, key: str = None) -> str:
        """
        Retrieves a secret from the provider.
        path: The path or identifier for the secret.
        key: Optional specific key within a secret object/map.
        """
        pass

class EnvSecretProvider(SecretProvider):
    def get_secret(self, path: str, key: str = None) -> str:
        # For Env provider, path IS the env var name. Key is ignored.
        val = os.environ.get(path)
        if val is None:
            logger.warning(f"Secret not found in environment: {path}")
            return ""
        return val

class MockVaultProvider(SecretProvider):
    def get_secret(self, path: str, key: str = None) -> str:
        # Simulates a vault. Path is like "secret-engine/path/to/secret"
        # Since this is a mock, we'll just return a dummy string based on path.
        logger.info(f"MockVault: Retrieving {path}")
        return f"mock-vault-value-for-{path}"
