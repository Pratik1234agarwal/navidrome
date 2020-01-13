package db_sql

import (
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cloudsonic/sonic-server/domain"
)

type MediaFile struct {
	ID          string    `orm:"pk;column(id)"`
	Path        string    ``
	Title       string    `orm:"index"`
	Album       string    ``
	Artist      string    ``
	ArtistID    string    `orm:"column(artist_id)"`
	AlbumArtist string    ``
	AlbumID     string    `orm:"column(album_id);index"`
	HasCoverArt bool      ``
	TrackNumber int       ``
	DiscNumber  int       ``
	Year        int       ``
	Size        string    ``
	Suffix      string    ``
	Duration    int       ``
	BitRate     int       ``
	Genre       string    ``
	Compilation bool      ``
	PlayCount   int       ``
	PlayDate    time.Time `orm:"null"`
	Rating      int       ``
	Starred     bool      `orm:"index"`
	StarredAt   time.Time `orm:"null"`
	CreatedAt   time.Time `orm:"null"`
	UpdatedAt   time.Time `orm:"null"`
}

type mediaFileRepository struct {
	sqlRepository
}

func NewMediaFileRepository() domain.MediaFileRepository {
	r := &mediaFileRepository{}
	r.tableName = "media_file"
	return r
}

func (r *mediaFileRepository) Put(m *domain.MediaFile) error {
	tm := MediaFile(*m)
	return r.put(m.ID, &tm)
}

func (r *mediaFileRepository) Get(id string) (*domain.MediaFile, error) {
	tm := MediaFile{ID: id}
	err := Db().Read(&tm)
	if err == orm.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	a := domain.MediaFile(tm)
	return &a, nil
}

func (r *mediaFileRepository) toMediaFiles(all []MediaFile) domain.MediaFiles {
	result := make(domain.MediaFiles, len(all))
	for i, m := range all {
		result[i] = domain.MediaFile(m)
	}
	return result
}

func (r *mediaFileRepository) FindByAlbum(albumId string) (domain.MediaFiles, error) {
	var mfs []MediaFile
	_, err := r.newQuery(Db()).Filter("album_id", albumId).OrderBy("disc_number", "track_number").All(&mfs)
	if err != nil {
		return nil, err
	}
	return r.toMediaFiles(mfs), nil
}

func (r *mediaFileRepository) GetStarred(options ...domain.QueryOptions) (domain.MediaFiles, error) {
	var starred []MediaFile
	_, err := r.newQuery(Db(), options...).Filter("starred", true).All(&starred)
	if err != nil {
		return nil, err
	}
	return r.toMediaFiles(starred), nil
}

func (r *mediaFileRepository) PurgeInactive(activeList domain.MediaFiles) ([]string, error) {
	return r.purgeInactive(activeList, func(item interface{}) string {
		return item.(domain.MediaFile).ID
	})
}

func (r *mediaFileRepository) Search(q string, offset int, size int) (domain.MediaFiles, error) {
	parts := strings.Split(q, " ")
	if len(parts) == 0 {
		return nil, nil
	}
	qs := r.newQuery(Db(), domain.QueryOptions{Offset: offset, Size: size})
	cond := orm.NewCondition()
	for _, part := range parts {
		c := orm.NewCondition()
		cond = cond.AndCond(c.Or("title__istartswith", part).Or("title__icontains", " "+part))
	}
	qs = qs.SetCond(cond).OrderBy("-rating", "-starred", "-play_count")
	var results []MediaFile
	_, err := qs.All(&results)
	if err != nil {
		return nil, err
	}
	return r.toMediaFiles(results), nil
}

var _ domain.MediaFileRepository = (*mediaFileRepository)(nil)
var _ = domain.MediaFile(MediaFile{})