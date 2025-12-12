#[cfg(test)]
mod pak_repository_test {
    use super::*;
    use tempfile::{tempdir, TempDir};
    use rand::Rng;
    use application::repository::blob_repository::BlobRepository;
    use persistence::repository::pak_repository::PakRepository;

    fn random_bytes(size: usize) -> Vec<u8> {
        let mut rng = rand::rng();
        (0..size).map(|_| rng.random::<u8>()).collect()
    }

    struct TestRepo {
        _temp_dir: TempDir, // удерживаем TempDir живым
        repo: PakRepository,
    }

    fn mock_pak_repository() -> TestRepo {
        let temp_dir = tempdir().unwrap();
        let key = *b"01234567890123456789012345678901";
        let repo = PakRepository::new(temp_dir.path(), key).unwrap();
        TestRepo { _temp_dir: temp_dir, repo }
    }

    #[tokio::test]
    async fn test_put_get() {
        let tr = mock_pak_repository();
        let key = "next_key";
        let data = b"some test data";

        tr.repo.put(key, data.as_ref()).await.unwrap();
        let retrieved_data = tr.repo.get(key).await.unwrap();
        assert_eq!(retrieved_data, data);
    }

    #[tokio::test]
    async fn test_get_metadata() {
        let tr = mock_pak_repository();
        let key = "meta_key";
        let data = b"metadata test data";

        tr.repo.put(key, data).await.unwrap();
        let meta = tr.repo.get_metadata(key).await.unwrap();

        assert_eq!(meta.key, key);
        assert_eq!(meta.size_original, data.len() as u64);
        assert!(meta.size_compressed > 0);
    }

    #[tokio::test]
    async fn test_exists_and_list() {
        let tr = mock_pak_repository();
        let key1 = "exist_key1";
        let key2 = "exist_key2";
        let data = b"existence test data";

        tr.repo.put(key1, data).await.unwrap();
        tr.repo.put(key2, data).await.unwrap();

        let list = tr.repo.list(None).await.unwrap();
        assert!(list.contains_key(key1));
        assert!(list.contains_key(key2));
    }

    #[tokio::test]
    async fn test_wrong_key_fails() {
        let temp_dir = tempdir().unwrap();
        let correct_key = *b"01234567890123456789012345678901";
        let wrong_key = *b"11111111111111111111111111111111";

        let repo = PakRepository::new(temp_dir.path(), correct_key).unwrap();
        let key = "secure";
        let data = b"secret";

        repo.put(key, data).await.unwrap();

        let repo_wrong = PakRepository::new(temp_dir.path(), wrong_key).unwrap();
        let result = repo_wrong.get(key).await;

        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_delete_not_supported() {
        let tr = mock_pak_repository();
        let result = tr.repo.delete("any").await;
        assert!(result.is_err());
    }
}
