FROM alpine

WORKDIR /app

COPY --from=build_location:develop /app/cmd/location/main ./app

CMD ["/app/app"]
