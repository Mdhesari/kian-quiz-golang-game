package mongorbac

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DB) CreateRole(ctx context.Context, role entity.Role) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	_, err := d.cli.Conn().Collection("roles").InsertOne(ctx, role)
	if err != nil {

		return err
	}

	return nil
}

func (d *DB) GetRoles(ctx context.Context) ([]entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	cur, err := d.cli.Conn().Collection("roles").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	roles := []entity.Role{}
	cur.All(ctx, &roles)

	return roles, nil
}

func (d *DB) GetRole(ctx context.Context, name string) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	res := d.cli.Conn().Collection("roles").FindOne(ctx, bson.M{
		"name": name,
	})
	if res.Err() != nil {

		return nil, res.Err()
	}

	var role entity.Role
	res.Decode(&role)

	return &role, nil
}

func (d *DB) GetPermissionIds(ctx context.Context, perms ...string) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	cur, err := d.cli.Conn().Collection("permissions").Find(ctx, bson.M{
		"name": bson.M{"$in": perms},
	})
	if err != nil {
		return nil, err
	}

	permissions := []entity.Permission{}
	if err := cur.All(ctx, permissions); err != nil {
		return nil, err
	}

	ids := []primitive.ObjectID{}
	for _, perm := range permissions {
		ids = append(ids, perm.ID)
	}

	return ids, nil
}

func (d *DB) HasPermissions(ctx context.Context, roleID primitive.ObjectID, permissionIDs ...primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	cur, err := d.cli.Conn().Collection("access_controls").Find(ctx, bson.M{
		"role_id":       roleID,
		"permission_id": bson.M{"$in": permissionIDs},
	})
	if err != nil {

		return false, err
	}

	accesses := []entity.Access{}
	if err := cur.All(ctx, accesses); err != nil {

		return false, err
	}

	//TODO: check accesses
	return false, nil
}
