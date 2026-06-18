from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8", extra="ignore")

    db_host: str = "127.0.0.1"
    db_port: int = 5432
    db_user: str = "your_db_user"
    db_password: str = "your_db_password"
    db_name: str = "pharmacy_erp"

    llm_base_url: str = "http://127.0.0.1:1234/v1"
    llm_api_key: str = "replace_with_your_llm_api_key"
    llm_model: str = "your_llm_model"

    service_host: str = "0.0.0.0"
    service_port: int = 9080

    service_version: str = "1.0.0"

    @property
    def dsn(self) -> str:
        return (
            f"postgresql://{self.db_user}:{self.db_password}"
            f"@{self.db_host}:{self.db_port}/{self.db_name}"
        )


settings = Settings()
