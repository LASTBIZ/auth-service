package biz

import (
	"auth-service/internal/biz"
	"auth-service/internal/mocks/mrepo"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HashUseCase", func() {
	var hashUseCase *biz.HashUseCase
	var mockHashRepo *mrepo.MockHashRepo

	BeforeEach(func() {
		mockHashRepo = mrepo.NewMockHashRepo(ctl)
		hashUseCase = biz.NewHashUseCase(mockHashRepo, nil)
	})

	It("Create", func() {
		hash := &biz.Hash{UserID: 1, Hash: "hash"}
		mockHashRepo.EXPECT().CreateHash(ctx, gomock.Any()).Return(hash, nil)
		hash1, err := hashUseCase.Create(ctx, hash)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(err).ToNot(HaveOccurred())
		Ω(hash1.UserID).To(Equal(uint32(1)))
		Ω(hash1.Hash).To(Equal("hash"))
	})

})
