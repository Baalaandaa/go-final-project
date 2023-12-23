FROM alpine

WORKDIR /app

COPY --from=build_driver:develop /app/cmd/driver/main ./app

CMD ["/app/app"]
