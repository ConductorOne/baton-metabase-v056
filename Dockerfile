FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-metabase-v056"]
COPY baton-metabase-v056 /