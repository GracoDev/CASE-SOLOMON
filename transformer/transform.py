import time

def main():
    print("Transformer service started")
    print("Service: transformer")
    print("Status: running")
    print("Message: Hello from Transformer (Data Transformation Service)")
    
    # Mantém o serviço rodando
    while True:
        time.sleep(10)
        print(f"Transformer heartbeat: {time.strftime('%Y-%m-%d %H:%M:%S')}")

if __name__ == '__main__':
    main()



