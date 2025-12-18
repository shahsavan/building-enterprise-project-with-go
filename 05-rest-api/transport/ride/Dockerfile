# ---------- Build stage ----------
    FROM golang:1.25-alpine AS builder

    # Set working directory inside the container
    WORKDIR /app
    
    # Install build tools required by Makefile
    RUN apk add --no-cache make
    
    # Copy the rest of the source code
    COPY . .
    
    # Build the application using the Makefile, placing the binary in /out
    RUN BIN_DIR=/out make build
    
    # ---------- Runtime stage ----------
    FROM gcr.io/distroless/base-debian12
    
    # Set working directory for the runtime container
    WORKDIR /app
    
    # Copy the compiled binary from the builder stage
    COPY --from=builder /out/ride /app/ride
    
    # Expose the HTTP port used by the web application
    EXPOSE 8080
    
    # Run the application
    ENTRYPOINT ["/app/ride"]
    
