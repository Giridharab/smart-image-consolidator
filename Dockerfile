FROM python:3.12
WORKDIR /app
RUN apt-get update && apt-get install -y curl git build-essential
RUN pip install flask requests numpy pandas
COPY . .
EXPOSE 5000
CMD ["python", "app.py"]
