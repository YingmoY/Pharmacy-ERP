$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

if (-not (Test-Path ".venv\Scripts\python.exe")) {
    Write-Host "Creating virtual environment..."
    python -m venv .venv
}

Write-Host "Installing / verifying dependencies..."
.\.venv\Scripts\pip install -r requirements.txt --trusted-host mirrors.ustc.edu.cn --trusted-host mirrors.tuna.tsinghua.edu.cn -q

Write-Host "Starting AI service..."
.\.venv\Scripts\python main.py
