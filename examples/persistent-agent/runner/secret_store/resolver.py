import logging
from urllib.parse import urlparse, parse_qs
from .provider import SecretProvider, EnvSecretProvider, MockVaultProvider

logger = logging.getLogger("gantral-runner")

class SecretResolver:
    def __init__(self):
        self.providers = {
            "env": EnvSecretProvider(),
            "vault": MockVaultProvider()
        }
    
    def resolve(self, value: str) -> str:
        """
        Resolves a value if it is a secret reference (gantral+secret://).
        Otherwise returns the value as-is.
        """
        if not isinstance(value, str) or not value.startswith("gantral+secret://"):
            return value

        try:
            # Parse URI: gantral+secret://provider/path?key=k
            # urlparse expects scheme to be just the protocol, but here it handles Custom+Scheme fine for parsing body.
            parsed = urlparse(value)
            
            provider_name = parsed.netloc
            path = parsed.path.lstrip('/') # Remove leading slash
            query = parse_qs(parsed.query)
            key = query.get('key', [None])[0]

            logger.info(f"Resolving secret: {value}") # Log reference only!

            provider = self.providers.get(provider_name)
            if not provider:
                logger.error(f"Unknown secret provider: {provider_name}")
                return "" # Fail safe

            return provider.get_secret(path, key)

        except Exception as e:
            logger.error(f"Failed to resolve secret URI {value}: {e}")
            return ""
